package serv

import (
	"github.com/ichaly/yugong/core/data"
	"go.uber.org/fx"
	"time"
)

type SpiderParams struct {
	fx.In
	Spiders []Spider `group:"spider"`
}

type Spider interface {
	Name() data.Platform
	GetAuthor(author *data.Author) error
	GetVideos(openId string, aid string, max *time.Time, min *time.Time) error
}
