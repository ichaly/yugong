package douyin

import (
	"errors"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/ichaly/yugong/core/base"
	"github.com/ichaly/yugong/core/data"
	"github.com/ichaly/yugong/core/serv"
	"github.com/ichaly/yugong/core/util"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
	"net/url"
	"regexp"
	"strings"
	"time"
)

type Douyin struct {
	db     *gorm.DB
	config *base.Config
	script *Script
}

func NewDouyin(d *gorm.DB, c *base.Config, s *Script) *Douyin {
	return &Douyin{db: d, config: c, script: s}
}

func (my Douyin) Name() data.Platform {
	return data.DouYin
}

func (my Douyin) GetAuthor(author *data.Author) error {
	req := serv.NewFetch(my.config).NoRedirect().UseProxy()
	reg := regexp.MustCompile(`[a-z]+://[\S]+`)
	uri := reg.FindAllString(author.Url, -1)[0]
	res, err := req.Get(uri)
	if err != nil {
		return err
	}
	if res.StatusCode() != 302 {
		return errors.New("not 302")
	}
	location := res.Header().Values("location")[0]
	regNew := regexp.MustCompile(`(?:sec_uid=)[a-z,A-Z，0-9, _, -]+`)
	sec_uid := strings.Replace(regNew.FindAllString(location, -1)[0], "sec_uid=", "", 1)
	str := fmt.Sprintf("https://www.iesdouyin.com/web/api/v2/user/info/?sec_uid=%s", sec_uid)
	var body string
	err = retry.Do(func() error {
		res, err := req.Get(str)
		if err != nil {
			return err
		}
		body = res.String()
		if body == "" {
			return errors.New("get user info body is empty")
		}
		return nil
	})
	if err != nil {
		return err
	}
	info := gjson.Get(body, "user_info")
	if !info.Exists() {
		return errors.New("get user info body is empty")
	}
	author.From = data.DouYin
	author.OpenId = sec_uid
	author.Fid = info.Get("uid").String()
	author.Nickname = info.Get("nickname").String()
	author.Signature = info.Get("signature").String()
	author.Avatar = info.Get("avatar_medium.url_list.0").String()
	my.db.Save(author)
	return nil
}

func (my Douyin) GetVideos(openId string, aid string, max *time.Time, min *time.Time) error {
	params := url.Values{"sec_user_id": []string{openId}, "count": []string{"50"}, "aid": []string{"6383"}}
	if min != nil {
		params.Add("min_cursor", fmt.Sprintf("%d", min.UnixNano()/1e6))
	}
	if max != nil {
		params.Add("max_cursor", fmt.Sprintf("%d", max.UnixNano()/1e6))
	}
	uri, _ := url.Parse("https://www.douyin.com/aweme/v1/web/aweme/post/")
	uri.RawQuery = params.Encode()
	req := serv.NewFetch(my.config).UseProxy().SetHeaders(map[string]string{
		"referer": "https://www.douyin.com/",
	}).SetCookies(map[string]string{
		"passport_csrf_token": "c8b96614139f50d240232221b574cacb",
		"ttwid":               "1%7CHQXlIa0A7vFQ2Je4UliR5vOoYX6tSdv24RZqMfNUaFg%7C1680791252%7Cb04d0bdf0e025c3135156fd23fdc730398da79188b3b458b2a051ca326dc962f",
	})
	str := util.JoinString(uri.String(), "&X-Bogus=", my.script.Sign(uri.RawQuery, req.Agent))
	var body string
	err := retry.Do(func() error {
		res, err := req.Get(str)
		if err != nil {
			return err
		}
		body = res.String()
		if body == "" {
			return errors.New("get videos body is empty")
		}
		return nil
	})
	if err != nil {
		return err
	}
	list := gjson.Get(body, "aweme_list").Array()
	var videos []data.Video
	for _, r := range list {
		uid := r.Get("author.uid").String()
		vid := r.Get("aweme_id").String()
		title := r.Get("desc").String()
		video := r.Get("video.play_addr.url_list.0").String()
		cover := r.Get("video.cover.url_list|@reverse|0").String()
		width := r.Get("video.width").Int()
		height := r.Get("video.height").Int()
		createTime := r.Get("create_time").Int()
		uploadTime := time.Now()
		v := data.Video{
			From: data.DouYin, Vid: vid, Url: video, Title: title, Cover: cover, Width: width, Height: height,
			Fid: uid, Aid: aid, UploadAt: util.TimePtr(uploadTime), SourceAt: time.UnixMilli(createTime * 1000),
		}
		videos = append(videos, v)
	}
	if len(videos) > 0 {
		my.db.Save(videos)
		if max != nil {
			max = util.TimePtr(time.UnixMilli(gjson.Get(body, "max_cursor").Int() * 1000))
		}
		if min != nil {
			min = util.TimePtr(time.UnixMilli(gjson.Get(body, "min_cursor").Int() * 1000))
		}
		//err := my.GetVideos(openId, aid, min, max)
		//if err != nil {
		//	return err
		//}
	}
	return nil
}
