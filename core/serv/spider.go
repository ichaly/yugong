package serv

import (
	"github.com/ichaly/yugong/core/data"
	"go.uber.org/fx"
	"time"
)

type SpiderGroup struct {
	fx.In
	Spiders []Spider `group:"spider"`
}

type Spider interface {
	Name() data.Platform
	GetAuthor(author *data.Author) error
	GetVideos(fid, aid string, cursor, finish *string, start *time.Time, total, count int) error
}
