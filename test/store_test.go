package test

import (
	gocache "github.com/eko/gocache/lib/v4/cache"
	"github.com/ichaly/yugong/core/base"
	"testing"
)

var store *gocache.Cache[string]

func init() {
	store, _ = base.NewStore(config)
}

func TestStore(t *testing.T) {
	t.Logf("store init:%+v", store)
}
