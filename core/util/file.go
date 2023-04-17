package util

import (
	"io"
	"os"
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
