package base

import (
	"context"
	"fmt"
	"github.com/allegro/bigcache/v3"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/store"
	sb "github.com/eko/gocache/store/bigcache/v4"
	sr "github.com/eko/gocache/store/redis/v4"
	"github.com/redis/go-redis/v9"
	"strings"
	"time"
)

func NewStore(c *Config) (*cache.Cache[string], error) {
	var s store.StoreInterface
	if strings.ToLower(c.Cache.Dialect) == "redis" {
		args := []interface{}{c.Cache.Host, c.Cache.Port}
		s = sr.NewRedis(redis.NewClient(&redis.Options{
			Addr:     fmt.Sprintf("%s:%d", args...),
			Username: c.Cache.Username,
			Password: c.Cache.Password,
		}))
	} else {
		client, err := bigcache.New(context.Background(), bigcache.DefaultConfig(5*time.Minute))
		if err != nil {
			return nil, err
		}
		s = sb.NewBigcache(client)
	}
	return cache.New[string](s), nil
}
