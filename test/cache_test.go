package test

import (
	"github.com/ichaly/yugong/core/base"
	"testing"
)

var cache base.Cache

func init() {
	cache = base.NewCache(store)
}

func TestCache(t *testing.T) {
	t.Logf("cache init:%+v", cache)
}
