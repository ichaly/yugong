package util

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/go-resty/resty/v2"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func WriteFile(source io.Reader, target string) error {
	_ = os.MkdirAll(filepath.Dir(target), 0777)
	file, err := os.OpenFile(target, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		_ = file.Close()
	}(file)
	_, err = io.Copy(file, source)
	if err != nil {
		return err
	}
	return nil
}

// UploadFile https://market.m.taobao.com/app/tb-zhibo-app/mcn-center/daren.html#/
func UploadFile(source, aid, padding string) error {
	client := resty.New()
	fileName := path.Base(source)
	fileBytes, err := os.ReadFile(source)
	if err != nil {
		return err
	}
	res, err := client.R().
		SetFileReader("file", fileName, bytes.NewReader(fileBytes)).
		SetFormData(map[string]string{
			"key":            fmt.Sprintf("daren/2215630453359/zip/%s%s", padding, fileName),
			"policy":         "eyJleHBpcmF0aW9uIjoiMjAyNS0wNC0xN1QxMzowOTozOC4zMjNaIiwiY29uZGl0aW9ucyI6W1siY29udGVudC1sZW5ndGgtcmFuZ2UiLDAsMjA5NzE1MjAwMF0sWyJzdGFydHMtd2l0aCIsIiRrZXkiLCJkYXJlbi8yMjE1NjMwNDUzMzU5L3ppcC8iXV19",
			"Signature":      "8uCPEIVsVa9o8cV4RXAUoScBX94=",
			"OSSAccessKeyId": "LTAI4G6jrdfWbfGFZJGNYnKN",
		}).Post("https://duke-daren-videos.oss-cn-beijing.aliyuncs.com/")
	if err != nil {
		return err
	}
	if find := strings.Contains(res.String(), "Error"); find {
		return errors.New(fmt.Sprintf("upload file failed.:%s=>%s", source, res.String()))
	}
	return nil
}

func DownloadFile(url string, target string) (err error) {
	client := resty.New().SetDoNotParseResponse(true)
	res, err := client.R().Get(url)
	if err != nil {
		return
	}
	err = WriteFile(res.RawBody(), target)
	if err != nil {
		return
	}
	return
}
