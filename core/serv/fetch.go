package serv

import (
	"fmt"
	"github.com/EDDYCJY/fake-useragent"
	"github.com/go-resty/resty/v2"
	"github.com/tidwall/gjson"
	"log"
)

type Fetch struct {
	Agent string
}

func NewFetch() *Fetch {
	authKey := "EF40DC7E"
	password := "8CEF0614705B"
	proxyServer := "123.54.54.175:23022"
	targetUrl := "https://ip.cn/api/index?ip=&type=0"
	proxyUrl := fmt.Sprintf("http://%s:%s@%s", authKey, password, proxyServer)
	//https://share.proxy.qg.net/pool?key=EF40DC7E&num=1&area=&isp=&format=json&seq=&pool=1
	client := resty.New().SetProxy(proxyUrl)
	f := Fetch{browser.Random()}
	res, err := client.R().SetHeader("user-Agent", f.Agent).Get(targetUrl)
	if err != nil {
		log.Println(err)
	}
	text := gjson.GetBytes(res.Body(), "ip").String()
	log.Println(text)
	return &f
}
