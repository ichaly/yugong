package test

import (
	"github.com/go-resty/resty/v2"
	"github.com/ichaly/yugong/core/serv"
	"net/url"
	"strings"
	"testing"
)

func TestHttp(t *testing.T) {
	s, err := serv.NewScript()
	if err != nil {
		return
	}
	client := resty.New()
	uri, _ := url.Parse("https://www.douyin.com/aweme/v1/web/aweme/post/")
	params := url.Values{
		"aid":         []string{"6383"},
		"sec_user_id": []string{"MS4wLjABAAAA69ZgRVFTFzxrD9LqFs3jCiZEbg1F7Ox8B4SbY5_Ver8"},
		"max_cursor":  []string{"1671719358000"},
		"count":       []string{"1"},
	}
	uri.RawQuery = params.Encode()
	res, err := client.R().
		SetHeader("user-agent", s.Agent()).
		SetHeader("referer", "https://www.douyin.com/").
		SetHeader("cookie", "passport_csrf_token=c8b96614139f50d240232221b574cacb;ttwid=1%7CHQXlIa0A7vFQ2Je4UliR5vOoYX6tSdv24RZqMfNUaFg%7C1680791252%7Cb04d0bdf0e025c3135156fd23fdc730398da79188b3b458b2a051ca326dc962f").
		Get(strings.Join([]string{uri.String(), "&X-Bogus=", s.Sign(uri.RawQuery)}, ""))
	if err != nil {
		panic(err)
	}
	println(res.String())
}
