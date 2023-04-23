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
	"strconv"
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

func (my Douyin) GetAuthor(url string) (map[string]string, error) {
	req := serv.NewFetch(my.config).NoRedirect().UseProxy()
	reg := regexp.MustCompile(`[a-z]+://[\S]+`)
	url = reg.FindAllString(url, -1)[0]
	res, err := req.Get(url)
	if err != nil {
		return nil, err
	}
	if res.StatusCode() != 302 {
		return nil, errors.New("not 302")
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
		return nil, err
	}
	info := gjson.Get(body, "user_info")
	return map[string]string{
		"openid":                    sec_uid,
		"uid":                       info.Get("uid").String(),
		"nickname":                  info.Get("nickname").String(),
		"signature":                 info.Get("signature").String(),
		"avatar":                    info.Get("avatar_larger.url_list.0").String(),
		"aweme_count":               info.Get("aweme_count").String(),
		"mplatform_followers_count": info.Get("mplatform_followers_count").String(),
	}, nil
}

func (my Douyin) GetVideos(openId string, did string, aid string, min int64, max int64) (int64, int64, error) {
	params := url.Values{
		"sec_user_id": []string{openId},
		"count":       []string{"30"},
		"aid":         []string{"6383"},
	}
	if min > 0 {
		params.Add("min_cursor", strconv.FormatInt(min, 10))
	}
	if max <= 0 {
		max = time.Now().UnixNano() / 1e6
	}
	params.Add("max_cursor", strconv.FormatInt(max, 10))
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
		return 0, 0, err
	}
	list := gjson.Get(body, "aweme_list").Array()
	var videos []*data.Video
	for _, r := range list {
		vid := r.Get("aweme_id").String()
		title := r.Get("desc").String()
		video := r.Get("video.play_addr.url_list.0").String()
		cover := r.Get("video.cover.url_list|@reverse|0").String()
		createTime := r.Get("create_time").Int()
		uploadTime := time.Now()
		v := data.Video{
			From: data.DouYin, Title: title, Url: video, Fid: did, Aid: aid, Cover: cover,
			Vid: vid, UploadAt: util.TimePtr(uploadTime), SourceAt: time.UnixMilli(createTime * 1000),
		}
		videos = append(videos, &v)
	}
	if len(videos) > 0 {
		my.db.Save(videos)
	}
	min = gjson.Get(body, "min_cursor").Int()
	max = gjson.Get(body, "max_cursor").Int()
	return min, max, nil
}
