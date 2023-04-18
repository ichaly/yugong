package util

import (
	"bytes"
	"fmt"
	"github.com/go-resty/resty/v2"
	"io"
	"os"
	"path"
	"path/filepath"
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
	key := fmt.Sprintf("daren/%s/zip/%s%s", aid, padding, fileName)
	println(key)
	res, err := client.R().
		SetFileReader("file", fileName, bytes.NewReader(fileBytes)).
		SetFormData(map[string]string{
			"key":            key,
			"policy":         "eyJleHBpcmF0aW9uIjoiMjAyNS0wNC0xN1QwOTo1NjoyMi41ODlaIiwiY29uZGl0aW9ucyI6W1siY29udGVudC1sZW5ndGgtcmFuZ2UiLDAsMjA5NzE1MjAwMF0sWyJzdGFydHMtd2l0aCIsIiRrZXkiLCJkYXJlbi82NTYyNjU2ODIvemlwLyJdXX0=",
			"Signature":      "Uvl/DXV5HWo2O251M4+vG7nHhkw=",
			"OSSAccessKeyId": "LTAI4G6jrdfWbfGFZJGNYnKN",
		}).Post("https://duke-daren-videos.oss-cn-beijing.aliyuncs.com/")
	if err != nil {
		return err
	}
	println(res.String())
	return nil
}
