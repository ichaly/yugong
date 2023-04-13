package serv

import (
	"errors"
	"fmt"
	"github.com/avast/retry-go"
	"github.com/kirinlabs/HttpRequest"
	"github.com/tidwall/gjson"
	"net/http"
	"regexp"
	"strings"
)

var (
	req *HttpRequest.Request
)

func init() {
	req = HttpRequest.NewRequest()
	req.SetHeaders(map[string]string{
		"User-Agent": "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/89.0.4389.114 Mobile Safari/537.36",
	})
	req.CheckRedirect(func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse /* 不进入重定向 */
	})
}

type Spider struct {
}

func NewSpider() *Spider {
	return &Spider{}
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

//https://github.com/PuerkitoBio/goquery
//https://github.com/gocolly/colly
//https://github.com/go-resty/resty
