package serv

import (
	"github.com/ichaly/yugong/core/data"
	"go.uber.org/fx"
)

type SpiderParams struct {
	fx.In
	Spiders []Spider `group:"spider"`
}

type Spider interface {
	Name() data.Platform
	GetAuthor(url string) (map[string]string, error)
	GetVideos(openId string, aid string, min int64, max int64) (int64, int64, error)
}
