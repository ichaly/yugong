package serv

import (
	"context"
	"github.com/go-co-op/gocron"
	"github.com/ichaly/yugong/core/base"
	"github.com/ichaly/yugong/core/data"
	"github.com/ichaly/yugong/core/util"
	"github.com/ichaly/yugong/zlog"
	"go.uber.org/fx"
	"gorm.io/gorm"
	"strconv"
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

func NewCrontab(l fx.Lifecycle, d *gorm.DB, c *base.Config, g SpiderGroup, q *Queue) *Crontab {
	timezone, _ := time.LoadLocation("Asia/Shanghai")
	scheduler := gocron.NewScheduler(timezone)
	scheduler.SingletonModeAll()
	crontab := &Crontab{
		db: d, config: c, queue: q,
		spiders:   make(map[data.Platform]Spider),
		scheduler: scheduler,
	}
	for _, s := range g.Spiders {
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
	tx.Where("disable = ?", false).Find(&authors)

	// add cron job
	for _, author := range authors {
		my.Watch(author)
	}

	// add sync job
	_, _ = my.scheduler.Tag(SYNC_FILES).Every(1).Hour().Do(my.syncFiles)

	my.scheduler.StartAsync()
}

func (my *Crontab) getVideos(authorId int64) {
	zlog.Debug("get videos", zlog.Int64("authorId", authorId))
	var author data.Author
	my.db.First(&author, authorId)
	var cursor, finish *string
	if author.From == data.DouYin {
		var maxTime time.Time
		row := my.db.Model(&data.Video{}).Select("max(source_at) as maxTime").Where("aid = ?", author.Aid).Row()
		_ = row.Scan(&maxTime)
		if !maxTime.IsZero() {
			finish = util.StringPtr(strconv.FormatInt(maxTime.UnixMilli(), 10))
		}
		cursor = util.StringPtr(strconv.FormatInt(time.Now().UnixMilli(), 10))
	} else if author.From == data.XiaoHongShu {
		row := my.db.Model(&data.Video{}).
			Select("vid").Where("aid = ?", author.Aid).
			Order("source_at desc").Limit(1)
		_ = row.Scan(&finish)
	}
	err := my.spiders[author.From].GetVideos(
		author.OpenId, author.Aid, cursor, finish, author.Start, author.Total, 0,
	)
	if err != nil {
		return
	}
}

func (my *Crontab) syncFiles() {
	var videos []data.Video
	my.db.Where("state", 0).Order("source_at asc").Find(&videos)
	for _, v := range videos {
		my.queue.Push(NewTask(my.config.Workspace, my.db, v))
	}
}
