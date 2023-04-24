package rest

import (
	"github.com/bytedance/sonic"
	"github.com/go-chi/chi/v5"
	"github.com/ichaly/yugong/core/base"
	"github.com/ichaly/yugong/core/data"
	"github.com/ichaly/yugong/core/serv"
	"gorm.io/gorm"
	"io"
	"net/http"
)

type DouyinApi struct {
	db      *gorm.DB
	render  *base.Render
	crontab *serv.Crontab
}

func NewDouyinApi(db *gorm.DB, rd *base.Render, c *serv.Crontab) base.Plugin {
	return &DouyinApi{db: db, render: rd, crontab: c}
}

func (my *DouyinApi) Name() string {
	return "DouyinApi"
}

func (my *DouyinApi) Protected() bool {
	return false
}

func (my *DouyinApi) Init(r chi.Router) {
	r.Route("/douyin", func(r chi.Router) {
		r.Get("/once", my.startHandler)
		r.Post("/save", my.saveHandler)
	})
}

func (my *DouyinApi) startHandler(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Query().Get("tag")
	my.crontab.Once(tag)
	_ = my.render.JSON(w, base.OK.WithMessage("操作成功"))
}

func (my *DouyinApi) saveHandler(w http.ResponseWriter, r *http.Request) {
	var u data.Author
	bty, err := io.ReadAll(r.Body)
	if err != nil {
		_ = my.render.JSON(w, base.ERROR.WithError(err), base.WithCode(http.StatusBadRequest))
		return
	}
	err = sonic.Unmarshal(bty, &u)
	if err != nil {
		_ = my.render.JSON(w, base.ERROR.WithError(err), base.WithCode(http.StatusBadRequest))
		return
	}
	if u.Aid == "" || u.Url == "" {
		_ = my.render.JSON(w, base.ERROR.WithMessage("参数aid或url不能为空"))
		return
	}
	err = my.crontab.GetSpider(data.DouYin).GetAuthor(&u)
	if err != nil {
		_ = my.render.JSON(w, base.ERROR.WithError(err))
		return
	}
	my.crontab.Watch(u)
	_ = my.render.JSON(w, base.OK.WithData(u))
}
