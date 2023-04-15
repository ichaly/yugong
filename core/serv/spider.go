package serv

import (
	_ "embed"
	"errors"
	"fmt"
	"github.com/avast/retry-go"
	_ "github.com/ddliu/motto/underscore"
	"github.com/dop251/goja"
	"github.com/kirinlabs/HttpRequest"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

var (
	req *HttpRequest.Request
)

func init() {
	req = HttpRequest.NewRequest()
	req.SetHeaders(map[string]string{
		"user-agent": agent,
		"referer":    "https://www.douyin.com/user/MS4wLjABAAAA69ZgRVFTFzxrD9LqFs3jCiZEbg1F7Ox8B4SbY5_Ver8",
		//"cookie":     "passport_csrf_token=c8b96614139f50d240232221b574cacb;ttwid=1%7CHQXlIa0A7vFQ2Je4UliR5vOoYX6tSdv24RZqMfNUaFg%7C1680791252%7Cb04d0bdf0e025c3135156fd23fdc730398da79188b3b458b2a051ca326dc962",
	})
	req.CheckRedirect(func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse /* 不进入重定向 */
	})
}

type Spider struct {
	db *gorm.DB
}

func NewSpider(d *gorm.DB) *Spider {
	return &Spider{d}
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

func (my *Spider) GetVideos(did string, aid string, maxCursor string, count int) error {
	vm := goja.New()
	_, err := vm.RunString(bogus)
	if err != nil {
		return err
	}
	var fn func(string, string) string
	err = vm.ExportTo(vm.Get("sign"), &fn)
	if err != nil {
		return err
	}
	uri, _ := url.Parse("https://www.douyin.com/aweme/v1/web/aweme/post/")
	params := url.Values{
		"device_platform":             []string{"webapp"},
		"aid":                         []string{"6383"},
		"channel":                     []string{"channel_pc_web"},
		"sec_user_id":                 []string{"MS4wLjABAAAA69ZgRVFTFzxrD9LqFs3jCiZEbg1F7Ox8B4SbY5_Ver8"},
		"max_cursor":                  []string{"1671719358000"},
		"locate_query":                []string{"false"},
		"show_live_replay_strategy":   []string{"1"},
		"count":                       []string{"10"},
		"publish_video_strategy_type": []string{"2"},
		"pc_client_type":              []string{"1"},
		"version_code":                []string{"170400"},
		"version_name":                []string{"17.4.0"},
		"cookie_enabled":              []string{"true"},
		"screen_width":                []string{"1440"},
		"screen_height":               []string{"900"},
		"browser_language":            []string{"zh-CN"},
		"browser_platform":            []string{"MacIntel"},
		"browser_name":                []string{"Edge"},
		"browser_version":             []string{"112.0.1722.39"},
		"browser_online":              []string{"true"},
		"engine_name":                 []string{"Blink"},
		"engine_version":              []string{"112.0.0.0"},
		"os_name":                     []string{"Mac OS"},
		"os_version":                  []string{"10.15.7"},
		"cpu_core_num":                []string{"8"},
		"device_memory":               []string{"8"},
		"platform":                    []string{"PC"},
		"downlink":                    []string{"10"},
		"effective_type":              []string{"4g"},
		"round_trip_time":             []string{"250"},
		"webid":                       []string{"7218943418981516855"},
		"msToken":                     []string{"5dz4s1lo5sG4J-7044vNPrSAxjh7S759AhN5XI0PFC6pa3Xsem9RdaHGfKMfJXUTKq5IvBZ88YTk4b58TVVxS1B5cjKT3Ozn4gnU6PYklAKrDYVGPo8UHHM8qQ=="},
		"X-Bogus":                     []string{"DFSzsdVE5dUANyEutVGFa6J22JAr"},
	}
	uri.RawQuery = params.Encode()

	t := uri.String() + "?"
	for s := range params {
		t = t + s + "=" + params.Get(s) + "&"
	}
	t = t[0 : len(t)-1]
	t = "https://www.douyin.com/aweme/v1/web/aweme/post/?device_platform=webapp&aid=6383&channel=channel_pc_web&sec_user_id=MS4wLjABAAAA69ZgRVFTFzxrD9LqFs3jCiZEbg1F7Ox8B4SbY5_Ver8&max_cursor=1663921152000&locate_query=false&show_live_replay_strategy=1&count=10&publish_video_strategy_type=2&pc_client_type=1&version_code=170400&version_name=17.4.0&cookie_enabled=true&screen_width=1440&screen_height=900&browser_language=zh-CN&browser_platform=MacIntel&browser_name=Edge&browser_version=112.0.1722.39&browser_online=true&engine_name=Blink&engine_version=112.0.0.0&os_name=Mac+OS&os_version=10.15.7&cpu_core_num=8&device_memory=8&platform=PC&downlink=10&effective_type=4g&round_trip_time=250&webid=7218943418981516855&msToken=5dz4s1lo5sG4J-7044vNPrSAxjh7S759AhN5XI0PFC6pa3Xsem9RdaHGfKMfJXUTKq5IvBZ88YTk4b58TVVxS1B5cjKT3Ozn4gnU6PYklAKrDYVGPo8UHHM8qQ==&X-Bogus=DFSzsdVE5dUANyEutVGFa6J22JAr"
	//https://www.douyin.com/aweme/v1/web/aweme/post/?X-Bogus=DFSzsdVE5dUANyEutVGFa6J22JAr&aid=6383&browser_language=zh-CN&browser_name=Edge&browser_online=true&browser_platform=MacIntel&browser_version=112.0.1722.39&channel=channel_pc_web&cookie_enabled=true&count=10&cpu_core_num=8&device_memory=8&device_platform=webapp&downlink=10&effective_type=4g&engine_name=Blink&engine_version=112.0.0.0&locate_query=false&max_cursor=1671719358000&msToken=5dz4s1lo5sG4J-7044vNPrSAxjh7S759AhN5XI0PFC6pa3Xsem9RdaHGfKMfJXUTKq5IvBZ88YTk4b58TVVxS1B5cjKT3Ozn4gnU6PYklAKrDYVGPo8UHHM8qQ%3D%3D&os_name=Mac+OS&os_version=10.15.7&pc_client_type=1&platform=PC&publish_video_strategy_type=2&round_trip_time=250&screen_height=900&screen_width=1440&sec_user_id=MS4wLjABAAAA69ZgRVFTFzxrD9LqFs3jCiZEbg1F7Ox8B4SbY5_Ver8&show_live_replay_strategy=1&version_code=170400&version_name=17.4.0&webid=7218943418981516855?browser_name=Edge&engine_name=Blink&cpu_core_num=8&device_memory=8&round_trip_time=250&show_live_replay_strategy=1&X-Bogus=DFSzsdVE5dUANyEutVGFa6J22JAr&browser_online=true&screen_width=1440&browser_language=zh-CN&os_version=10.15.7&platform=PC&effective_type=4g&cookie_enabled=true&browser_platform=MacIntel&engine_version=112.0.0.0&os_name=Mac OS&version_name=17.4.0&locate_query=false&count=10&publish_video_strategy_type=2&browser_version=112.0.1722.39&downlink=10&webid=7218943418981516855&max_cursor=1671719358000&aid=6383&pc_client_type=1&version_code=170400&device_platform=webapp&sec_user_id=MS4wLjABAAAA69ZgRVFTFzxrD9LqFs3jCiZEbg1F7Ox8B4SbY5_Ver8&channel=channel_pc_web&msToken=5dz4s1lo5sG4J-7044vNPrSAxjh7S759AhN5XI0PFC6pa3Xsem9RdaHGfKMfJXUTKq5IvBZ88YTk4b58TVVxS1B5cjKT3Ozn4gnU6PYklAKrDYVGPo8UHHM8qQ==&screen_height=900
	//https://www.douyin.com/aweme/v1/web/aweme/post/?device_platform=webapp&aid=6383&channel=channel_pc_web&sec_user_id=MS4wLjABAAAA69ZgRVFTFzxrD9LqFs3jCiZEbg1F7Ox8B4SbY5_Ver8&max_cursor=1663921152000&locate_query=false&show_live_replay_strategy=1&count=10&publish_video_strategy_type=2&pc_client_type=1&version_code=170400&version_name=17.4.0&cookie_enabled=true&screen_width=1440&screen_height=900&browser_language=zh-CN&browser_platform=MacIntel&browser_name=Edge&browser_version=112.0.1722.39&browser_online=true&engine_name=Blink&engine_version=112.0.0.0&os_name=Mac+OS&os_version=10.15.7&cpu_core_num=8&device_memory=8&platform=PC&downlink=10&effective_type=4g&round_trip_time=250&webid=7218943418981516855&msToken=5dz4s1lo5sG4J-7044vNPrSAxjh7S759AhN5XI0PFC6pa3Xsem9RdaHGfKMfJXUTKq5IvBZ88YTk4b58TVVxS1B5cjKT3Ozn4gnU6PYklAKrDYVGPo8UHHM8qQ==&X-Bogus=DFSzsdVE5dUANyEutVGFa6J22JAr
	//p, _ := url.Parse(t)
	//e := fn(p.RawQuery, agent)
	//println(e)

	//enc := fn(uri.RawQuery, agent)
	//u := fmt.Sprintf("%s&X-Bogus=%s", uri.String(), enc)

	req = HttpRequest.NewRequest()
	req.SetHeaders(map[string]string{
		"user-agent":  agent,
		"user-cookie": "passport_csrf_token=c8b96614139f50d240232221b574cacb;ttwid=1%7CHQXlIa0A7vFQ2Je4UliR5vOoYX6tSdv24RZqMfNUaFg%7C1680791252%7Cb04d0bdf0e025c3135156fd23fdc730398da79188b3b458b2a051ca326dc962f",
		"referer":     "https://www.douyin.com/user/MS4wLjABAAAA69ZgRVFTFzxrD9LqFs3jCiZEbg1F7Ox8B4SbY5_Ver8",
	})
	//req.SetCookies(map[string]string{
	//	"passport_csrf_token": "c8b96614139f50d240232221b574cacb",
	//	"ttwid":               "1%7CHQXlIa0A7vFQ2Je4UliR5vOoYX6tSdv24RZqMfNUaFg%7C1680791252%7Cb04d0bdf0e025c3135156fd23fdc730398da79188b3b458b2a051ca326dc962f",
	//})
	//req.Proxy(func(request *http.Request) (*url.URL, error) {
	//	return url.Parse("http://127.0.0.1:8888")
	//})
	resp, err := req.Get(t)
	defer resp.Close()
	if err != nil {
		return err
	}
	body, err := resp.Body()
	if err != nil {
		return err
	}
	list := gjson.GetBytes(body, "aweme_list").Array()
	size := len(list)
	println(size)
	//if size > 0 {
	//	for i := 0; i < size; i++ {
	//		title := gjson.Get(string(body), fmt.Sprintf("aweme_list.%d.desc", i)).String()
	//		uri := gjson.Get(string(body), fmt.Sprintf("aweme_list.%d.video.play_addr.url_list.0", i)).String()
	//		v := &data.Video{
	//			Title: title,
	//			Url:   uri,
	//			Did:   did,
	//			Aid:   aid,
	//		}
	//		my.db.Save(v)
	//	}
	//	count += size
	//	if count < 30 {
	//		return my.GetVideos(did, aid, endTime, count)
	//	}
	//}
	return nil
}

//https://github.com/PuerkitoBio/goquery
//https://github.com/gocolly/colly
//https://github.com/go-resty/resty
//https://github.com/guonaihong/gout
