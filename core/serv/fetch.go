package serv

import (
	"github.com/EDDYCJY/fake-useragent"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-resty/resty/v2"
	"log"
)

type Fetch struct {
	agent string
}

func NewFetch() *Fetch {
	client := resty.New().SetDoNotParseResponse(true)
	f := Fetch{browser.Random()}
	res, err := client.R().
		SetHeader("User-Agent", browser.Random()).
		Get("https://2023.ip138.com")
	if err != nil {
		panic(err)
	}
	doc, err := goquery.NewDocumentFromReader(res.RawBody())
	if err != nil {
		panic(err)
	}
	text := doc.Find("body > p:nth-child(1) > a:nth-child(1)").Text()
	log.Println(text)
	return &f
}
