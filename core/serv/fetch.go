package serv

import (
	"fmt"
	"github.com/EDDYCJY/fake-useragent"
	"github.com/avast/retry-go"
	"github.com/go-resty/resty/v2"
	"github.com/ichaly/yugong/core/base"
	"github.com/tidwall/gjson"
	"net/http"
	"net/url"
)

type Fetch struct {
	Agent string

	cong    *base.Config
	client  *resty.Client
	cookies []*http.Cookie
	files   map[string]string
	params  map[string]string
	headers map[string]string
}

type noRedirectPolicy struct{}

func (my noRedirectPolicy) Apply(req *http.Request, via []*http.Request) error {
	return http.ErrUseLastResponse
}

func NewFetch(c *base.Config) *Fetch {
	agent := browser.Random()
	f := Fetch{
		Agent: agent,
		cong:  c, client: resty.New(),
		files:  make(map[string]string),
		params: make(map[string]string),
		headers: map[string]string{
			"user-agent": agent,
		},
	}
	return &f
}

func (my *Fetch) setProxy() error {
	var proxy string
	err := retry.Do(func() error {
		client := resty.New()
		params := url.Values{
			"num":    []string{"1"},
			"pool":   []string{"1"},
			"format": []string{"json"},
			"key":    []string{my.cong.Proxy.Username},
		}
		uri, err := url.Parse(my.cong.Proxy.Host)
		if err != nil {
			return err
		}
		uri.RawQuery = params.Encode()
		res, err := client.R().SetHeader("user-agent", my.Agent).Get(uri.String())
		if err != nil {
			return err
		}
		proxy = gjson.GetBytes(res.Body(), "data.0.server").String()
		return nil
	})
	proxyUrl := fmt.Sprintf("http://%s:%s@%s", my.cong.Proxy.Username, my.cong.Proxy.Password, proxy)
	my.client.SetProxy(proxyUrl)
	return err
}

func (my *Fetch) UseProxy() *Fetch {
	_ = my.setProxy()
	return my
}

func (my *Fetch) NoRedirect() *Fetch {
	/* 不进入重定向 */
	my.client.SetRedirectPolicy(noRedirectPolicy{})
	return my
}

func (my *Fetch) SetHeaders(headers map[string]string) *Fetch {
	for k, v := range headers {
		my.headers[k] = v
	}
	return my
}

func (my *Fetch) SetParams(params map[string]string) *Fetch {
	for k, v := range params {
		my.params[k] = v
	}
	return my
}

func (my *Fetch) SetCookies(cookies map[string]string) *Fetch {
	for k, v := range cookies {
		my.cookies = append(my.cookies, &http.Cookie{Name: k, Value: v})
	}
	return my
}

func (my *Fetch) SetFiles(files map[string]string) *Fetch {
	for k, v := range files {
		my.files[k] = v
	}
	return my
}

func (my *Fetch) Get(uri string) (*resty.Response, error) {
	return my.client.R().
		EnableTrace().
		SetCookies(my.cookies).
		SetHeaders(my.headers).
		SetQueryParams(my.params).
		Get(uri)
}

func (my *Fetch) Json(uri string) (*resty.Response, error) {
	return my.client.R().
		EnableTrace().
		SetCookies(my.cookies).
		SetHeaders(my.headers).
		SetBody(my.params).
		Get(uri)
}

func (my *Fetch) Form(uri string) (*resty.Response, error) {
	return my.client.R().
		EnableTrace().
		SetCookies(my.cookies).
		SetHeaders(my.headers).
		SetFormData(my.params).
		Get(uri)
}

func (my *Fetch) Upload(uri string) (*resty.Response, error) {
	return my.client.R().
		EnableTrace().
		SetCookies(my.cookies).
		SetHeaders(my.headers).
		SetFiles(my.files).
		SetFormData(my.params).
		Get(uri)
}

func (my *Fetch) Download(uri, desc string) (*resty.Response, error) {
	return my.client.R().
		EnableTrace().
		SetCookies(my.cookies).
		SetHeaders(my.headers).
		SetQueryParams(my.params).
		SetOutput(desc).
		Get(uri)
}
