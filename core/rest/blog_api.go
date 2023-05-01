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

type BlogApi struct {
	db      *gorm.DB
	render  *base.Render
	crontab *serv.Crontab
}

func NewBlogApi(db *gorm.DB, rd *base.Render, c *serv.Crontab) base.Plugin {
	return &BlogApi{db: db, render: rd, crontab: c}
}

func (my *BlogApi) Name() string {
	return "BlogApi"
}

func (my *BlogApi) Protected() bool {
	return false
}

func (my *BlogApi) Init(r chi.Router) {
	r.Route("/blog", func(r chi.Router) {
		r.Get("/once", my.startHandler)
		r.Post("/save", my.saveHandler)
	})
}

func (my *BlogApi) startHandler(w http.ResponseWriter, r *http.Request) {
	tag := r.URL.Query().Get("tag")
	my.crontab.Once(tag)
	_ = my.render.JSON(w, base.OK.WithMessage("操作成功"))
}

func (my *BlogApi) saveHandler(w http.ResponseWriter, r *http.Request) {
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
	if u.Aid == "" || u.Url == "" || u.From == "" {
		_ = my.render.JSON(w, base.ERROR.WithMessage("参数aid,url,from不能为空"))
		return
	}
	spider := my.crontab.GetSpider(u.From)
	if spider == nil {
		_ = my.render.JSON(w, base.ERROR.WithMessage("不支持的平台"))
		return
	}
	err = spider.GetAuthor(&u)
	if err != nil {
		_ = my.render.JSON(w, base.ERROR.WithError(err))
		return
	}
	my.crontab.Watch(u)
	_ = my.render.JSON(w, base.OK.WithData(u))
}
