package serv

import (
	"github.com/go-resty/resty/v2"
	"net/url"
	"strings"
	"testing"
)

func TestHttp(t *testing.T) {
	s, err := NewScript()
	if err != nil {
		return
	}
	client := resty.New()
	uri, _ := url.Parse("https://www.douyin.com/aweme/v1/web/aweme/post/")
	params := url.Values{
		"device_platform":             []string{"webapp"},
		"aid":                         []string{"6383"},
		"channel":                     []string{"channel_pc_web"},
		"sec_user_id":                 []string{"MS4wLjABAAAA69ZgRVFTFzxrD9LqFs3jCiZEbg1F7Ox8B4SbY5_Ver8"},
		"max_cursor":                  []string{"1671719358000"},
		"locate_query":                []string{"false"},
		"show_live_replay_strategy":   []string{"1"},
		"count":                       []string{"1"},
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
		//"X-Bogus":                     []string{"DFSzsdVE5dUANyEutVGFa6J22JAr"},
	}
	uri.RawQuery = params.Encode()
	res, err := client.R().
		SetHeader("user-agent", agent).
		SetHeader("referer", "https://www.douyin.com/").
		SetHeader("cookie", "passport_csrf_token=c8b96614139f50d240232221b574cacb;ttwid=1%7CHQXlIa0A7vFQ2Je4UliR5vOoYX6tSdv24RZqMfNUaFg%7C1680791252%7Cb04d0bdf0e025c3135156fd23fdc730398da79188b3b458b2a051ca326dc962f").
		Get(strings.Join([]string{uri.String(), "&X-Bogus=", s.Sign(uri.RawQuery)}, ""))
	if err != nil {
		panic(err)
	}
	println(res.String())
}
