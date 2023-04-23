package serv

import (
	"context"
	"github.com/go-co-op/gocron"
	"github.com/ichaly/yugong/core/base"
	"github.com/ichaly/yugong/core/data"
	"github.com/ichaly/yugong/core/util"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"time"
)

type Crontab struct {
	db        *gorm.DB
	queue     *Queue
	config    *base.Config
	scheduler *gocron.Scheduler
	spiders   map[data.Platform]Spider
}

func NewCrontab(l fx.Lifecycle, d *gorm.DB, c *base.Config, s SpiderParams, q *Queue) *Crontab {
	timezone, _ := time.LoadLocation("Asia/Shanghai")
	crontab := &Crontab{
		db: d, config: c, queue: q,
		spiders:   make(map[data.Platform]Spider),
		scheduler: gocron.NewScheduler(timezone),
	}
	for _, s := range s.Spiders {
		crontab.spiders[s.Support()] = s
	}
	l.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			crontab.Start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			crontab.Stop()
			return nil
		},
	})
	return crontab
}

func (my *Crontab) Stop() {
	my.scheduler.Stop()
}

func (my *Crontab) Start() {
	// get authors by condition
	var authors []*data.Author
	tx := my.db
	if my.config.Condition != nil {
		tx = tx.Where(my.config.Condition.Query, my.config.Condition.Values)
	}
	tx.Find(&authors)

	// group by cron
	group := make(map[string][]*data.Author)
	for _, author := range authors {
		if v, ok := group[author.Cron]; ok {
			group[author.Cron] = append(v, author)
		} else {
			group[author.Cron] = []*data.Author{author}
		}
	}

	// add cron job
	for k, _ := range group {
		if k == "" {
			_, _ = my.scheduler.Every(1).Day().At("00:00").Do(func() {
				for _, author := range group[k] {
					my.GetVideos(author)
				}
			})
		} else {
			_, _ = my.scheduler.Cron(k).Do(func() {
				for _, author := range group[k] {
					my.GetVideos(author)
				}
			})
		}
	}

	// add sync job
	_, _ = my.scheduler.Every(1).Hour().Do(my.SyncFiles)

	my.scheduler.StartAsync()
}

func (my *Crontab) GetSpider(p data.Platform) Spider {
	return my.spiders[p]
}

func (my *Crontab) GetVideos(a *data.Author) {
	go func(a *data.Author) {
		var oldMin, oldMax int64
		if a.MaxTime != nil {
			oldMax = a.MaxTime.UnixNano() / 1e6
		}
		if a.MinTime != nil {
			oldMin = a.MinTime.UnixNano() / 1e6
		}
		newMin, newMax, err := my.spiders[a.From].GetVideos(
			a.OpenId, a.Fid, a.Aid, oldMin, oldMax,
		)
		if err != nil {
			return
		}
		if newMin < oldMin || oldMin == 0 {
			a.MinTime = util.TimePtr(time.UnixMilli(newMin))
		}
		if newMax > oldMax {
			a.MaxTime = util.TimePtr(time.UnixMilli(newMax))
		}
		my.db.Save(&a)
	}(a)
}

func (my *Crontab) SyncFiles() {
	var videos []*data.Video
	my.db.Where("state", 0).Find(&videos)
	for _, v := range videos {
		my.queue.Push(NewTask(my.config.Workspace, my.db, *v))
	}
}
