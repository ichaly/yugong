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
	GetAuthor(author *data.Author) error
	GetVideos(openId string, aid string, min int64, max int64) error
}
