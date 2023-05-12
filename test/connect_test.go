package test

import (
	"github.com/ichaly/yugong/core/base"
	"gorm.io/gorm"
	"testing"
)

var connect *gorm.DB

func init() {
	connect, _ = base.NewConnect(cache, config)
}

func TestConnect(t *testing.T) {
	t.Logf("connect init:%+v", connect)
}
