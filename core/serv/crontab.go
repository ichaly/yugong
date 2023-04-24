package serv

import (
	"context"
	"github.com/go-co-op/gocron"
	"github.com/ichaly/yugong/core/base"
	"github.com/ichaly/yugong/core/data"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"time"
)

const (
	GET_VIDEOS = "getVideos"
	SYNC_FILES = "syncFiles"
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

func (my *Crontab) Watch(a data.Author) {
	if a.Cron == "" {
		_, _ = my.scheduler.Tag(GET_VIDEOS).Every(1).Day().At("00:00").Do(my.getVideos, a.Id)
	} else {
		_, _ = my.scheduler.Tag(GET_VIDEOS).Cron(a.Cron).Do(my.getVideos, a.Id)
	}
}

func (my *Crontab) Stop() {
	my.scheduler.Stop()
}

func (my *Crontab) Once(tag string) {
	if tag == "" {
		my.scheduler.RunAll()
	} else if tag == GET_VIDEOS {
		_ = my.scheduler.RunByTag(GET_VIDEOS)
	} else if tag == SYNC_FILES {
		_ = my.scheduler.RunByTag(SYNC_FILES)
	}
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
	_, _ = my.scheduler.Tag(SYNC_FILES).Every(1).Hour().Do(my.syncFiles)

	my.scheduler.StartAsync()
}

func (my *Crontab) getVideos(authorId int64) {
	var author data.Author
	my.db.First(&author, authorId)
	var max *time.Time
	var oldMin, oldMax int64
	row := my.db.Model(&data.Video{}).Select("max(source_at) as max").Where("aid = ?", author.Aid).Row()
	_ = row.Scan(&max)
	if max != nil {
		oldMin = max.UnixNano() / 1e6
	} else {
		oldMax = time.Now().UnixNano() / 1e6
	}
	err := my.spiders[author.From].GetVideos(author.OpenId, author.Aid, oldMax, oldMin)
	if err != nil {
		return
	}
}

func (my *Crontab) syncFiles() {
	var videos []data.Video
	my.db.Where("state", 0).Order("source_at desc").Find(&videos)
	for _, v := range videos {
		my.queue.Push(NewTask(my.config.Workspace, my.db, v))
	}
}
