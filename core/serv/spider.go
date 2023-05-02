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
	GetVideos(openId, aid string, more bool, cursor *string, start *time.Time, total, count int) error
}
