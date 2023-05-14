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
	"math/rand"
	"net/url"
	"strings"
	"time"
)

const (
	SESSION_KEY = "SESSION_KEY"
)

type BodyParseError struct {
	code    int64
	message string
}

func (err BodyParseError) Error() string {
	return err.message
}

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

func (my XiaoHongShu) GetVideos(aid, fid string, cursor, finish *string, start *time.Time, total, count int) error {
	if finish == nil && start == nil && total == 0 {
		return nil
	}
	params := url.Values{"num": []string{"50"}, "user_id": []string{fid}, "cursor": []string{""}}
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
		if r.Get("type").String() != "video" {
			continue
		}
		vid := r.Get("note_id").String()
		cover := fmt.Sprintf("https://sns-img-bd.xhscdn.com/%s", r.Get("cover.trace_id").String())
		title := r.Get("display_title").String()
		sticky := r.Get("interact_info.sticky").Bool()
		v := data.Video{
			From: data.XiaoHongShu, Vid: vid, Fid: fid, Aid: aid, Sticky: sticky,
			Title: title, Cover: cover, UploadAt: util.TimePtr(time.Now()),
		}
		if finish != nil {
			if strings.Compare(vid, *finish) <= 0 {
				break
			}
		} else if start != nil {
			if v.SourceAt != nil && start.UnixMilli() >= v.SourceAt.UnixMilli() {
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
		err := my.GetVideos(fid, aid, cursor, finish, start, total, count)
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

func (my XiaoHongShu) GetDetail(v *data.Video) error {
	//开始获取详情
	params := map[string]string{"source_note_id": v.Vid}
	uri, _ := url.Parse("https://edith.xiaohongshu.com/api/sns/web/v1/feed")
	token := my.script.Sign(uri.Path, params)
	req := serv.NewFetch(my.config).UseProxy().SetHeaders(map[string]string{
		"referer": "https://www.xiaohongshu.com/",
		"x-t":     token["X-t"],
		"x-s":     token["X-s"],
	}).SetParams(params)
	var body string
	zlog.Debug("开始请求详情",
		zlog.String("vid", v.Vid),
		zlog.String("platform", string(v.From)),
	)
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
	}, retry.OnRetry(func(n uint, err error) {
		s := rand.Intn(5)
		time.Sleep(time.Second * time.Duration(int(n+1)*5+s))
		req = req.UseProxy()
	}))
	if err != nil {
		return err
	}
	zlog.Info("结束请求详情",
		zlog.String("vid", v.Vid),
		zlog.String("body", body),
		zlog.String("platform", string(v.From)),
	)
	detail := gjson.Get(body, "data.items.0.note_card")
	v.SourceAt = util.TimePtr(time.UnixMilli(detail.Get("time").Int()))
	v.Url = fmt.Sprintf("http://sns-video-bd.xhscdn.com/%s", detail.Get("video.consumer.origin_video_key").String())
	return nil
}

func (my XiaoHongShu) session() (string, error) {
	ses, err := my.cache.Get(context.Background(), SESSION_KEY)
	if err == nil && ses != "" {
		return ses, nil
	}
	uri, _ := url.Parse("https://edith.xiaohongshu.com/api/sns/web/v1/login/activate")
	a1 := "187a433ad934uk8iu6eapr38z41whjm51j84it9ja30000145194"
	header := my.script.Header(a1, uri, nil)
	req := serv.NewFetch(my.config).UseProxy().SetHeaders(header).SetCookies(map[string]string{"a1": a1})
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
	code := gjson.Get(body, "code").Int()
	if code == -100 {
		//_ = my.cache.Delete(context.Background(), SESSION_KEY)
		return BodyParseError{code, "登录已过期"}
	} else if code == 300015 {
		return BodyParseError{code, "浏览器异常，请尝试关闭/卸载风险插件或重启试试"}
	} else if code == 300012 {
		return BodyParseError{code, "网络连接异常，请检查网络设置或重启试试"}
	} else if code != 0 {
		return BodyParseError{code, "result code is not 0"}
	}
	return nil
}
