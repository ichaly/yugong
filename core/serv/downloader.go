package serv

import (
	"github.com/go-resty/resty/v2"
	"io"
	"os"
	"path"
)

type DownloaderOption func(*Downloader)

type Downloader struct {
	output     string
	maxThread  int
	retryTimes int
}

func WithOutput(output string) DownloaderOption {
	return func(d *Downloader) {
		if output != "" {
			d.output = output
		}
	}
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

func NewDownloader(opts ...DownloaderOption) (Downloader, error) {
	d := Downloader{maxThread: 1, retryTimes: 3}

	tmp, err := os.MkdirTemp("", "YuGong*.tmp")
	if err != nil {
		return d, nil
	}

	//default setting
	d.output = tmp

	//custom setting
	for _, o := range opts {
		o(&d)
	}

	_ = os.Mkdir(d.output, 0777)

	return d, nil
}

func (my *Downloader) Download(url string, name string) (*os.File, error) {
	client := resty.New().SetDoNotParseResponse(true)
	res, err := client.R().Get(url)
	if err != nil {
		return nil, err
	}
	return my.WriteFile(res.RawBody(), name)
}

func (my *Downloader) WriteFile(src io.Reader, name string) (*os.File, error) {
	file, err := os.OpenFile(path.Join(my.output, name), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	defer func() {
		_ = file.Close()
	}()
	_, err = io.Copy(file, src)
	if err != nil {
		return nil, err
	}
	return file, nil
}
