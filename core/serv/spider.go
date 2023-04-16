package serv

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/go-resty/resty/v2"
	"github.com/ichaly/yugong/core/base"
	"github.com/ichaly/yugong/core/data"
	"github.com/ichaly/yugong/core/util"
	"github.com/kirinlabs/HttpRequest"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	req *HttpRequest.Request
)

func init() {
	req = HttpRequest.NewRequest()
	req.SetHeaders(map[string]string{
		"user-agent": agent,
		"referer":    "https://www.douyin.com/user/MS4wLjABAAAA69ZgRVFTFzxrD9LqFs3jCiZEbg1F7Ox8B4SbY5_Ver8",
	})
	req.CheckRedirect(func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse /* 不进入重定向 */
	})
}

type Spider struct {
	db     *gorm.DB
	script *Script
	queue  *Queue
	config *base.Config
}

func NewSpider(d *gorm.DB, s *Script, q *Queue, c *base.Config) *Spider {
	return &Spider{d, s, q, c}
}

func (my *Spider) GetUserInfo(url string) (map[string]string, error) {
	reg := regexp.MustCompile(`[a-z]+://[\S]+`)
	url = reg.FindAllString(url, -1)[0]
	resp, err := req.Get(url)
	defer resp.Close()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 302 {
		return nil, err
	}
	location := resp.Headers().Values("location")[0]
	regNew := regexp.MustCompile(`(?:sec_uid=)[a-z,A-Z，0-9, _, -]+`)
	sec_uid := strings.Replace(regNew.FindAllString(location, -1)[0], "sec_uid=", "", 1)
	var body []byte
	err = retry.Do(func() error {
		res, err := req.Get(fmt.Sprintf("https://www.iesdouyin.com/web/api/v2/user/info/?sec_uid=%s", sec_uid))
		defer res.Close()
		if err != nil {
			return err
		}
		body, err = res.Body()
		if err != nil {
			return err
		}
		if string(body) == "" {
			return errors.New("body is empty")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	info := gjson.ParseBytes(body).Get("user_info")
	return map[string]string{
		"did":                       sec_uid,
		"nickname":                  info.Get("nickname").String(),
		"signature":                 info.Get("signature").String(),
		"avatar":                    info.Get("avatar_larger.url_list.0").String(),
		"aweme_count":               info.Get("aweme_count").String(),
		"mplatform_followers_count": info.Get("mplatform_followers_count").String(),
	}, nil
}

func (my *Spider) GetVideos(did string, aid string, min int64) (int64, error) {
	client := resty.New()
	uri, _ := url.Parse("https://www.douyin.com/aweme/v1/web/aweme/post/")
	params := url.Values{
		"sec_user_id": []string{did},
		"count":       []string{"31"},
		"aid":         []string{"6383"},
		"max_cursor":  []string{strconv.FormatInt(time.Now().UnixNano()/1e6, 10)},
	}
	if min > 0 {
		params.Add("min_cursor", strconv.FormatInt(min, 10))
	}
	uri.RawQuery = params.Encode()
	res, err := client.R().
		SetHeader("user-agent", agent).
		SetHeader("referer", "https://www.douyin.com/").
		SetHeader("cookie", "passport_csrf_token=c8b96614139f50d240232221b574cacb;ttwid=1%7CHQXlIa0A7vFQ2Je4UliR5vOoYX6tSdv24RZqMfNUaFg%7C1680791252%7Cb04d0bdf0e025c3135156fd23fdc730398da79188b3b458b2a051ca326dc962f").
		Get(strings.Join([]string{uri.String(), "&X-Bogus=", my.script.Sign(uri.RawQuery)}, ""))
	if err != nil {
		return 0, err
	}
	body := res.String()
	list := gjson.Get(body, "aweme_list").Array()
	size := len(list)
	if size > 0 {
		for i := 0; i < size; i++ {
			vid := gjson.Get(body, fmt.Sprintf("aweme_list.%d.aweme_id", i)).String()
			title := gjson.Get(body, fmt.Sprintf("aweme_list.%d.desc", i)).String()
			video := gjson.Get(body, fmt.Sprintf("aweme_list.%d.video.play_addr.url_list.0", i)).String()
			cover := gjson.Get(body, fmt.Sprintf("aweme_list.%d.video.cover.url_list|@reverse|0", i)).String()
			source := gjson.Get(body, fmt.Sprintf("aweme_list.%d.create_time", i)).Int()
			v := &data.Video{Title: title, Url: video, Did: did, Aid: aid, Cover: cover, SourceAt: time.UnixMilli(source * 1000)}
			my.db.Save(v)
			my.queue.Add(func() {
				workspace := my.config.Workspace
				d, _ := NewDownloader(WithOutput(workspace))

				titleFile, err := d.WriteFile(strings.NewReader(v.Title), fmt.Sprintf("t_%s.txt", vid))
				if err != nil {
					return
				}
				defer os.Remove(titleFile.Name())
				videoFile, err := d.Download(v.Url, fmt.Sprintf("v_%s.mp4", vid))
				if err != nil {
					return
				}
				defer os.Remove(videoFile.Name())
				coverFile, err := d.Download(v.Cover, fmt.Sprintf("c_%s.jpg", vid))
				if err != nil {
					return
				}
				defer os.Remove(coverFile.Name())
				err = util.Compress(util.Join(workspace, vid, ".zip"), titleFile, videoFile, coverFile)
				if err != nil {
					return
				}
			})
		}
	}
	return gjson.Get(body, "min_cursor").Int(), nil
}

//https://github.com/PuerkitoBio/goquery
//https://github.com/gocolly/colly
//https://github.com/guonaihong/gout