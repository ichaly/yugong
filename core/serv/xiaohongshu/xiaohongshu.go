package xiaohongshu

import (
	"context"
	"errors"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/ichaly/yugong/core/base"
	"github.com/ichaly/yugong/core/data"
	"github.com/ichaly/yugong/core/serv"
	"github.com/ichaly/yugong/core/util"
	"github.com/ichaly/yugong/zlog"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
	"net/url"
	"strings"
	"time"
)

const (
	SESSION_KEY = "SESSION_KEY"
)

type XiaoHongShu struct {
	db     *gorm.DB
	config *base.Config
	script *Script
	cache  *cache.Cache[string]
}

func NewXiaoHongShu(d *gorm.DB, c *base.Config, s *Script, e *cache.Cache[string]) *XiaoHongShu {
	return &XiaoHongShu{db: d, config: c, script: s, cache: e}
}

func (my XiaoHongShu) Name() data.Platform {
	return data.XiaoHongShu
}

func (my XiaoHongShu) GetAuthor(author *data.Author) error {
	openid := strings.SplitN(author.Url, "/", 6)[5]
	params := url.Values{"target_user_id": []string{openid}}
	uri, _ := url.Parse("https://edith.xiaohongshu.com/api/sns/web/v1/user/otherinfo")
	uri.RawQuery = params.Encode()
	req := serv.NewFetch(my.config).UseProxy().SetHeaders(map[string]string{
		"referer": "https://www.xiaohongshu.com/",
		"x-t":     "1682947892201",
		"x-s":     "0YsC1iavZ2w6O6M+slkkOiT+OYFp1laB0Y1Csidvs6M3",
	})
	var body string
	err := retry.Do(func() error {
		session, err := my.session()
		if err != nil {
			return err
		}
		req = req.SetCookies(map[string]string{"web_session": session})
		res, err := req.Get(uri.String())
		if err != nil {
			return err
		}
		body = res.String()
		return my.check(body)
	})
	if err != nil {
		return err
	}
	info := gjson.Get(body, "data.basic_info")
	if !info.Exists() {
		return errors.New("get user info body is empty")
	}
	author.From = data.XiaoHongShu
	author.OpenId = openid
	author.Nickname = info.Get("nickname").String()
	author.Signature = info.Get("desc").String()
	author.Avatar = info.Get("imageb").String()
	my.db.Save(author)
	return nil
}

func (my XiaoHongShu) GetVideos(openId, aid string, cursor, finish *string, start *time.Time, total, count int) error {
	if finish == nil && start == nil && total == 0 {
		return nil
	}
	params := url.Values{"num": []string{"50"}, "user_id": []string{openId}, "cursor": []string{""}}
	if finish == nil && total > 0 {
		params.Set("num", fmt.Sprintf("%d", util.Min(50, total-count)))
	}
	if cursor != nil {
		params.Set("cursor", fmt.Sprintf("%s", *cursor))
	}
	uri, _ := url.Parse("https://edith.xiaohongshu.com/api/sns/web/v1/user_posted")
	uri.RawQuery = params.Encode()
	token := my.script.Sign(fmt.Sprintf("%s?%s", uri.Path, uri.RawQuery), nil)
	req := serv.NewFetch(my.config).UseProxy().SetHeaders(map[string]string{
		"referer": "https://www.xiaohongshu.com/",
		"x-t":     token["X-t"],
		"x-s":     token["X-s"],
	})
	var body string
	zlog.Debug("开始请求", zlog.String("uri", uri.String()))
	err := retry.Do(func() error {
		session, err := my.session()
		if err != nil {
			return err
		}
		req = req.SetCookies(map[string]string{"web_session": session})
		res, err := req.Get(uri.String())
		if err != nil {
			return err
		}
		body = res.String()
		return my.check(body)
	})
	if err != nil {
		return err
	}
	zlog.Info("结束请求", zlog.String("uri", uri.String()))
	list := gjson.Get(body, "data.notes").Array()
	videos := make([]data.Video, 0)
	for i, r := range list {
		typ := r.Get("type").String()
		if typ != "video" {
			continue
		}
		isTop := r.Get("interact_info.sticky").Bool()
		// TODO: 置顶数据暂时忽略
		if isTop {
			continue
		}

		vid := r.Get("note_id").String()
		cover := fmt.Sprintf("https://sns-img-bd.xhscdn.com/%s", r.Get("cover.trace_id").String())
		width := r.Get("cover.width").Int()
		height := r.Get("cover.height").Int()
		title := r.Get("display_title").String()
		v := data.Video{
			From: data.XiaoHongShu, Vid: vid, Title: title, Cover: cover, Width: width,
			Height: height, Fid: openId, Aid: aid, UploadAt: util.TimePtr(time.Now()),
		}
		err := my.detail(&v)
		if err != nil {
			return err
		}
		if finish != nil {
			if strings.Compare(vid, *finish) <= 0 {
				break
			}
		} else if start != nil {
			if start.UnixMilli() >= v.SourceAt.UnixMilli() {
				// 到达了开始时间
				break
			}
		} else if total != -1 && count+i >= total {
			// 达到了同步数量
			break
		}
		videos = append([]data.Video{v}, videos...)
	}

	size := len(videos)
	if size > 0 {
		count = count + size
		cursor = util.StringPtr(gjson.Get(body, "data.cursor").String())
		err := my.GetVideos(openId, aid, cursor, finish, start, total, count)
		if err != nil {
			return err
		}
		err = my.db.Save(videos).Error
		if err != nil {
			return err
		}
	}
	return nil
}

func (my XiaoHongShu) detail(v *data.Video) error {
	params := map[string]string{"source_note_id": v.Vid}
	uri, _ := url.Parse("https://edith.xiaohongshu.com/api/sns/web/v1/feed")
	token := my.script.Sign(uri.Path, params)
	req := serv.NewFetch(my.config).UseProxy().SetHeaders(map[string]string{
		"referer": "https://www.xiaohongshu.com/",
		"x-t":     token["X-t"],
		"x-s":     token["X-s"],
	}).SetParams(params)
	var body string
	zlog.Debug("开始请求详情", zlog.String("vid", v.Vid))
	err := retry.Do(func() error {
		session, err := my.session()
		if err != nil {
			return err
		}
		req = req.SetCookies(map[string]string{"web_session": session})
		res, err := req.Json(uri.String())
		if err != nil {
			return err
		}
		body = res.String()
		return my.check(body)
	})
	if err != nil {
		return err
	}
	zlog.Info("结束请求详情", zlog.String("vid", v.Vid))
	detail := gjson.Get(body, "data.items.0.note_card")
	if !detail.IsObject() {
		return errors.New("detail is not object")
	}
	v.SourceAt = time.UnixMilli(detail.Get("time").Int())
	v.Url = fmt.Sprintf("http://sns-video-bd.xhscdn.com/%s", detail.Get("video.consumer.origin_video_key").String())
	return nil
}

func (my XiaoHongShu) session() (string, error) {
	ses, err := my.cache.Get(context.Background(), SESSION_KEY)
	if err == nil && ses != "" {
		return ses, nil
	}
	uri, _ := url.Parse("https://edith.xiaohongshu.com/api/sns/web/v1/login/activate")
	req := serv.NewFetch(my.config).UseProxy().SetHeaders(map[string]string{
		"sec-ch-ua":          "\"Microsoft Edge\";v=\"113\", \"Chromium\";v=\"113\", \"Not-A.Brand\";v=\"24\"",
		"x-t":                "1.6837E+12",
		"x-b3-traceid":       "d5b51f56afaff2cb",
		"sec-ch-ua-mobile":   "?0",
		"user-agent":         "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/113.0.0.0 Safari/537.36 Edg/113.0.1773.35",
		"content-type":       "application/json;charset=UTF-8",
		"accept":             "application/json, text/plain, */*",
		"x-s-common":         "2UQAPsHCPUIjqArjwjHjNsQhPsHCH0rjNsQhPaHCH0P1PUhIHjIj2eHjwjQ+GnPW/MPjNsQhPUHCHdYiqUMIGUM78nHjNsQh+sHCH0H1+shFHjIj2eLjwjHlweWI+eS0wBPAG/Q1GgWAqeplwB8UPBrUweQ08946Je8D8/pi40L9w/PIPeZI+eLM+eWMHjIj2eGjw0r9weP7Pec7P/rI+0rVHdW7H0ijnbSgg9pEadkYp9zMp/+y4LSxJ9Swprbk/r+t2fbg8opnaBl7nS+Q+DS187YQyg4knpYs4M+gLnSOyLiFGLY+4B+o/gzDPS8kanS7ynPUJBEjJbkVG9EwqBHU+BSOyLShanS7yn+oz0pjzASinD+Q+DSxnfM+PDMC/p4z2rMgL/pwzB+h/SzmPSkxLfT+2f47nnknyMkL//p+PSDA/p4+2bkrnfk+yDQxn/QnyDRLpfTyzMkknnkaybDU//zypFDM/FzyJLRoagY+yS83/Szz2SSgz/m+pB+7/Lz++pSLpfMOpBli/0Q82bkxzfYwzF8xngkbPSkLafl+PDkTnnkd+pkLafT8pr83np4zPDMg//b+zbbhnDzyybSCafSyzF83/fkByDELc/m+pMrUnp4+2pSTafY+pFFInfMwyDELpgS+zFpCnSzwJLRgn/Q8yfVl//QyJLRrz/pyyfzx/D4z2DRLn/myzbSh/Lzz4FEgpfM8pMQT/Lz34FhUagY+PS8x/pzb2rEC8BSw2SLFnfMtyMkr/fYyzbLlngkb2rMgLfl+ySQinDzQPMSTa/Qwpbb7ngk0+pkrag4OzbLl/Lz8PDRL//bypbLIngkayrECzgS8PDpC/p482bDUag4wprFFnDzyyFRLGAmwzBql/pzm4FMLyAp8yDkxngkDypkra/p+prrInnM+2DExnfS82DrM/SzQ2DMrzgYyJLFA/gkzPbSgnfT+zrkx/Sz+2DhUp/z+2SLUnfkmPDELnfSOzBlVnD4b4MST/gYwzBlV/LzwySkTLfkOzFM7nnknJLErp/++2fY3nnMp+LRLpgk82SQi/L+twaHVHdWhH0ija/PhqDYD87+xJ7mdag8Sq9zn494QcUT6aLpPJLQy+nLApd4G/B4BprShLA+jqg4bqD8S8gYDPBp3Jf+m2DMBnnEl4BYQyrkSL9z+2obl49zQ4DbApFQ0yo4c4ozdJ/c9aMpC2rSiPoPI/rTAydb7JdD7zbkQ4fRA2BQcydS04LbQyrTSzBr7q98DqBlc4g4+PS874d+64gmQc7pT/Sqha7kn4M+Qc94SyM8Fao4l4FzQzLMAanWI8gYPnnzQ2bPAag89qM4M4MH3+7kc2f+wq9zVP9pfLo4LanS9qM8c4rTOpdzlq7bFpFYVafpLJrq6ag8rzrSiqpzdpd4OadbF4LShzn+yL94SpdpFwLS94fpnpdzyanTVaFS3woQ1c/4AnpplPnQm+npxpdzCJMPMqM+0qS4C8epSLM40LrSh+nprqgzIaLpw8nzl4eSjnLz+4B8PzFShJ9pL4g4VJfFA8gYDzpmQyM4OanSHJgbAwrzQcFzltF8nJFSkG9SSLocMaL+yyezs/d+rqrTS2obFcLDAcg+x8A4SzBRyOaHVHdWEH0i7+ecEPePEw/cVHdWlPsHCP0DUKc==",
		"x-s":                "XYW_eyJzaWduU3ZuIjoiNTAiLCJzaWduVHlwZSI6IngxIiwiYXBwSWQiOiJ4aHMtcGMtd2ViIiwic2lnblZlcnNpb24iOiIxIiwicGF5bG9hZCI6IjZmM2MzMWQxMWQ5NDcxNTA2ZjRkMzgwYjVkZTM5M2Y3MWMxZDZjMjBjY2VjOTVkZGZiYjRhY2M4ZWI5ODZmOGIxMjVkNTQzYWE0MTcwMDM5YTVmOTlhN2YxZjFhNGFjYzE2ZTJlM2JmYjg5ZTJkYTFkYWQ2MWM1MDQxZDZhYzJiZGFkNjFjNTA0MWQ2YWMyYmJhMWM0ZmNjNTUyMGEzZTNmOWY2Yjk1M2ZmODE5ZjdjNGQzOTY0ZDYxMDQwNWVmYWRmMDkwN2IxM2VjMTExNzdiNzU4ZmJkZDNhZDU1YzExMWRlMjRhZDI3YmI2NTQwYzc5ZDIwODU1MDY2OTM1ZTU0YzRhNzEyY2EzMWYxY2IwNTM4ZDZkOTc0NDg1MTAwMTk5YjJjYzdiZDI5MTA0YmMzNjZiYzA5NTIzMDExZmM3MzQ0YWZkMDJjNTMzN2U4MzU2ZjA2NTZiODllZGEwYmMxNDllMDNjZmRjOGMwYjVmNDU3MzhkYmU5OTUzYzRhMCJ9",
		"sec-ch-ua-platform": "\"macOS\"",
		"origin":             "https://www.xiaohongshu.com",
		"sec-fetch-site":     "same-site",
		"sec-fetch-mode":     "cors",
		"sec-fetch-dest":     "empty",
		"referer":            "https://www.xiaohongshu.com/",
		"accept-encoding":    "gzip, deflate, br",
		"accept-language":    "zh-CN,zh;q=0.9,en;q=0.8,en-US;q=0.7",
	}).SetCookies(map[string]string{
		"acw_tc":            "5a9645e75fccbc9d6a87110a36d15e7c22b18a97a1d3e639461661dfe6d1dd0d",
		"xhsTrackerId":      "2fd9d3f2-85a1-4a82-98d3-dada8acc5c6b",
		"xhsTrackerId.sig":  "6Eh9OK7-pzw9jzGy3ZpBm1EixJarNLCWWUURJeW31sQ",
		"extra_exp_ids":     "yamcha_0327_clt,h5_1208_exp3,ques_clt2",
		"extra_exp_ids.sig": "-9P_FIY9nRpp4czlpi3JlPCL_zdr5ZMYd73Vy8sdzzY",
		"webBuild":          "2.4.4",
		"xsecappid":         "xhs-pc-web",
		"a1":                "188049c8c3a2nax3p5q8fr0a282cggol6de5hv56930000455485",
		"gid":               "yYY84jSYfYKYyYY84jSYSfh2q0JI0CqU2AYidKjF80JYJ2q8Skkl1x8884224Y28JqYjY82y",
		"gid.sign":          "/7dqcU70e3EnTailQnkW+HlfMvQ=",
	})
	var body string
	zlog.Debug("开始登陆:" + uri.String())
	res, err := req.Json(uri.String())
	if err != nil {
		return "", err
	}
	body = res.String()
	zlog.Info("结束登陆", zlog.String("body", body))
	ses = gjson.Get(body, "data.session").String()
	if ses == "" {
		return "", errors.New("session is empty")
	}
	err = my.cache.Set(context.Background(), SESSION_KEY, ses)
	if err != nil {
		return "", err
	}
	return ses, nil
}

func (my XiaoHongShu) check(body string) error {
	//gjson.Get(body, "msg").String() == "登录已过期" ||
	if gjson.Get(body, "code").Int() == -100 {
		//_ = my.cache.Delete(context.Background(), SESSION_KEY)
		return errors.New("登录已过期")
	}
	return nil
}
