package serv

import (
	"github.com/go-resty/resty/v2"
	"github.com/ichaly/yugong/core/util"
)

type DownloaderOption func(*Downloader)

type Downloader struct {
	maxThread  int
	retryTimes int
}

func WithMaxThread(maxThread int) DownloaderOption {
	return func(d *Downloader) {
		d.maxThread = maxThread
	}
}

func WithRetryTimes(retryTimes int) DownloaderOption {
	return func(d *Downloader) {
		d.retryTimes = retryTimes
	}
}

func NewDownloader(opts ...DownloaderOption) Downloader {
	d := Downloader{maxThread: 1, retryTimes: 3}

	//custom setting
	for _, o := range opts {
		o(&d)
	}

	return d
}

func (my *Downloader) Download(url string, target string) (err error) {
	client := resty.New().SetDoNotParseResponse(true)
	res, err := client.R().Get(url)
	if err != nil {
		return
	}
	err = util.WriteFile(res.RawBody(), target)
	if err != nil {
		return
	}
	return
}
