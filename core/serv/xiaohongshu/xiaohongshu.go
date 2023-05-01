package xiaohongshu

import (
	"errors"
	"github.com/avast/retry-go"
	"github.com/ichaly/yugong/core/base"
	"github.com/ichaly/yugong/core/data"
	"github.com/ichaly/yugong/core/serv"
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

func (my XiaoHongShu) GetVideos(openId, aid string, max, min, start *time.Time, total, count int) error {
	return nil
}
