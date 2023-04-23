package serv

import (
	"github.com/go-co-op/gocron"
	"time"
)

type Crontab struct {
	scheduler *gocron.Scheduler
}

func NewCrontab() *Crontab {
	timezone, _ := time.LoadLocation("Asia/Shanghai")
	return &Crontab{
		gocron.NewScheduler(timezone),
	}
}

func (my *Crontab) Run() {
	my.scheduler.StartAsync()
}
