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
	data := map[string]string{}
	uri, _ := url.Parse("https://edith.xiaohongshu.com/api/sns/web/v1/login/activate")
	a1 := "18805433dffvrl8b2hut9pmkcdwizrt067cr7g6i330000326172"
	header := my.script.Header(a1, uri, data)
	req := serv.NewFetch(my.config).UseProxy().SetHeaders(header).SetCookies(map[string]string{
		"a1":                a1,
		"xhsTrackerId":      "6b23733f-4cf5-41dd-9cca-49379fa85fdd",
		"xhsTrackerId.sig":  "Q8epTi0Bb_ZS-p7wWnXqmNMTm2BL19btcdna5GDdlSM",
		"webId":             "8f93932a0b9cd184743340efeb37f28e",
		"gid":               "yYY824qqiJK4yYY824qqfWkKiihF1YDJx79jS3qU6TSfEFq83uF98W888qJKyWJ8dfSD202Y",
		"gid.sign":          "t3eeEgSrn9qRxwSMjVtIpNI9cKs=",
		"xsecappid":         "xhs-pc-web",
		"web_session":       "040069b5511a2a147061ab2a69364bc558a1f3",
		"websectiga":        "6169c1e84f393779a5f7de7303038f3b47a78e47be716e7bec57ccce17d45f99",
		"sec_poison_id":     "eec1467b-f5c0-4155-9103-48ad9c266bf7",
		"extra_exp_ids":     "yamcha_0327_exp,h5_1208_exp3,ques_clt2",
		"extra_exp_ids.sig": "ETM51AFqVyLPOioG2x0qNaEzMLVwrEIN37uTpfkLqxc",
		"webBuild":          "2.4.5",
		"acw_tc":            "736474944edf09618c1a6293d17361df55ee01fad643d8c9eb0ebf95c5682002",
	}).SetParams(data)
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
