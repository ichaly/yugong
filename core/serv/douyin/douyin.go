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
	"log"
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

func (my Douyin) GetVideos(openId, aid string, cursor, finish *string, start *time.Time, total, count int) error {
	if finish == nil && start == nil && total == 0 {
		return nil
	}
	params := url.Values{"count": []string{"50"}, "sec_user_id": []string{openId}, "aid": []string{"6383"}}
	if finish == nil && total > 0 {
		params.Set("count", fmt.Sprintf("%d", util.Min(50, total-count)))
	}
	if cursor != nil {
		params.Add("max_cursor", *cursor)
	}
	if finish != nil {
		params.Add("min_cursor", *finish)
	}
	uri, _ := url.Parse("https://www.douyin.com/aweme/v1/web/aweme/post/")
	uri.RawQuery = params.Encode()
	req := serv.NewFetch(my.config).UseProxy().SetHeaders(map[string]string{
		"referer": "https://www.douyin.com/",
	}).SetCookies(map[string]string{
		"passport_csrf_token": "325b44eb177b269871d609f97649893e",
		"ttwid":               "1%7C8afBBEv3O1lekBnmOgwzAJSoHBy6kD7z_FahdsWeiLE%7C1680262762%7C37d3142c1a2b8e1eb10fb60cb2e88971b0088085bd36b1714a247063011cc77c",
	})
	str := util.JoinString(uri.String(), "&X-Bogus=", my.script.Sign(uri.RawQuery, req.Agent))
	var body string
	log.Println("开始请求:" + uri.String())
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
	log.Println("结束请求:" + uri.String())
	list := gjson.Get(body, "aweme_list").Array()
	videos := make([]data.Video, 0)
	for i, r := range list {
		isTop := r.Get("video.is_top").Int() == 1
		// TODO: 置顶数据暂时忽略
		if isTop {
			continue
		}
		uid := r.Get("author.uid").String()
		vid := r.Get("aweme_id").String()
		title := r.Get("desc").String()
		video := r.Get("video.play_addr.url_list.0").String()
		cover := r.Get("video.cover.url_list|@reverse|0").String()
		width := r.Get("video.width").Int()
		height := r.Get("video.height").Int()
		createTime := r.Get("create_time").Int() * 1000

		if finish != nil {
			// 到达了结束时间
			if createTime <= util.ParseLong(*finish) {
				break
			}
		} else if start != nil {
			// 到达了开始时间
			if start.UnixMilli() >= createTime {
				break
			}
		} else if total != -1 && count+i >= total {
			// 达到了同步数量
			break
		}

		v := data.Video{
			From: data.DouYin, Vid: vid, Url: video, Title: title, Cover: cover, Width: width, Height: height,
			Fid: uid, Aid: aid, UploadAt: util.TimePtr(time.Now()), SourceAt: time.UnixMilli(createTime),
		}
		videos = append([]data.Video{v}, videos...)
	}
	size := len(videos)
	if size > 0 {
		count = count + size
		cursor = util.StringPtr(strconv.FormatInt(gjson.Get(body, "max_cursor").Int(), 10))
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
