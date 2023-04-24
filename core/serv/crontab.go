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
	scheduler := gocron.NewScheduler(timezone)
	scheduler.SingletonModeAll()
	crontab := &Crontab{
		db: d, config: c, queue: q,
		spiders:   make(map[data.Platform]Spider),
		scheduler: scheduler,
	}
	for _, s := range s.Spiders {
		crontab.spiders[s.Name()] = s
	}
	l.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			crontab.start()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			crontab.Stop()
			return nil
		},
	})
	return crontab
}

func (my *Crontab) Watch(author data.Author) {
	if author.Cron == "" {
		_, _ = my.scheduler.Every(1).Day().At("00:00").Do(my.getVideos, author)
	} else {
		_, _ = my.scheduler.Cron(author.Cron).Do(my.getVideos, author)
	}
}

func (my *Crontab) Stop() {
	my.scheduler.Stop()
}

func (my *Crontab) Once() {
	my.scheduler.RunAll()
}

func (my *Crontab) GetSpider(p data.Platform) Spider {
	return my.spiders[p]
}

func (my *Crontab) start() {
	// get authors by condition
	tx := my.db
	var authors []data.Author
	if my.config.Condition != nil {
		tx = tx.Where(my.config.Condition.Query, my.config.Condition.Values)
	}
	tx.Find(&authors)

	// add cron job
	for _, author := range authors {
		my.Watch(author)
	}

	// add sync job
	_, _ = my.scheduler.Every(1).Hour().Do(my.syncFiles)

	my.scheduler.StartAsync()
}

func (my *Crontab) getVideos(a data.Author) {
	var oldMin, oldMax int64
	if a.MaxTime != nil {
		oldMax = a.MaxTime.UnixNano() / 1e6
	}
	if a.MinTime != nil {
		oldMin = a.MinTime.UnixNano() / 1e6
	}
	newMin, newMax, err := my.spiders[a.From].GetVideos(
		a.OpenId, a.Aid, oldMin, oldMax,
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
}

func (my *Crontab) syncFiles() {
	var videos []data.Video
	my.db.Where("state", 0).Find(&videos)
	for _, v := range videos {
		my.queue.Push(NewTask(my.config.Workspace, my.db, v))
	}
}
