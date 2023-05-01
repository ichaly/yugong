package xiaohongshu

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
	"strings"
	"time"
)

type XiaoHongShu struct {
	db     *gorm.DB
	config *base.Config
}

func NewXiaoHongShu(d *gorm.DB, c *base.Config) *XiaoHongShu {
	return &XiaoHongShu{db: d, config: c}
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
	}).SetCookies(map[string]string{
		"web_session": "040069b5511a2a147061d4f17a364b16fb5f6c",
	})
	var body string
	err := retry.Do(func() error {
		res, err := req.Get(uri.String())
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

func (my XiaoHongShu) GetVideos(openId, aid string, max, min *string, start *time.Time, total, count int) error {
	if max != nil && start == nil && total == 0 {
		return nil
	}
	page := "30"
	if max != nil {
		page = fmt.Sprintf("%d", util.Min(30, total-count))
	}
	params := url.Values{"num": []string{page}, "user_id": []string{openId}}
	if max != nil {
		params.Add("cursor", fmt.Sprintf("%s", *max))
	} else {
		params.Add("cursor", "")
	}
	req := serv.NewFetch(my.config).UseProxy().SetHeaders(map[string]string{
		"referer": "https://www.xiaohongshu.com/",
		"x-t":     "1682961327003",
		"x-s":     "OiOvsgqvslZJ1B5i0gvbO6aBsgTLsgwUZ2aB0gVUOlA3",
	}).SetCookies(map[string]string{
		"web_session": "040069b5511a2a147061d4f17a364b16fb5f6c",
	})
	uri, _ := url.Parse("https://edith.xiaohongshu.com/api/sns/web/v1/user_posted")
	uri.RawQuery = params.Encode()
	var body string
	err := retry.Do(func() error {
		res, err := req.Get(uri.String())
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
	list := gjson.Get(body, "data.notes").Array()
	videos := make([]data.Video, 0)
	for i, r := range list {
		typ := r.Get("type").String()
		if typ == "video" {
			continue
		}
		isTop := r.Get("interact_info.sticky").Bool()
		// 如果是置顶视频，要检测是否已经存在
		if isTop {
			continue
		}

		vid := r.Get("note_id").String()
		cover := r.Get("cover.url").String()
		width := r.Get("cover.width").Int()
		height := r.Get("cover.height").Int()
		title := r.Get("display_title").String()
		uploadTime := time.Now()
		v := data.Video{
			From: data.DouYin, Vid: vid, Title: title, Cover: cover, Width: width,
			Height: height, Fid: openId, Aid: aid, UploadAt: util.TimePtr(uploadTime),
		}
		err := my.detail(&v)
		if err != nil {
			return err
		}

		if max != nil {
			if start != nil && start.UnixMilli() >= v.SourceAt.UnixMilli() {
				// 到达了开始时间
				continue
			} else if total != -1 && count+i >= total {
				// 达到了同步数量
				break
			}
		}
		videos = append([]data.Video{v}, videos...)
	}
	size := len(videos)
	if size > 0 {
		count = count + size
		if max != nil {
			max = util.StringPtr(gjson.Get(body, "data.cursor").String())
			err := my.GetVideos(openId, aid, max, min, start, total, count)
			if err != nil {
				return err
			}
			err = my.db.Save(videos).Error
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (my XiaoHongShu) detail(v *data.Video) error {
	params := map[string]string{"source_note_id": v.Vid}
	req := serv.NewFetch(my.config).UseProxy().SetHeaders(map[string]string{
		"referer": "https://www.xiaohongshu.com/",
		"x-t":     "1682947892201",
		"x-s":     "0YsC1iavZ2w6O6M+slkkOiT+OYFp1laB0Y1Csidvs6M3",
	}).SetCookies(map[string]string{
		"web_session": "040069b5511a2a147061d4f17a364b16fb5f6c",
	}).SetParams(params)
	var body string
	err := retry.Do(func() error {
		res, err := req.Get("https://edith.xiaohongshu.com/api/sns/web/v1/feed")
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
	detail := gjson.Get(body, "data.items.0.note_card")
	if !detail.IsObject() {
		return errors.New("detail is not object")
	}
	v.SourceAt = time.UnixMilli(detail.Get("time").Int())
	v.Url = fmt.Sprintf("http://sns-video-bd.xhscdn.com/%s", detail.Get("video.consumer.origin_video_key").String())
	return nil
}
