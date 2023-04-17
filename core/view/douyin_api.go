package view

import (
	"github.com/bytedance/sonic"
	"github.com/go-chi/chi/v5"
	"github.com/ichaly/yugong/core/base"
	"github.com/ichaly/yugong/core/data"
	"github.com/ichaly/yugong/core/serv/douyin"
	"github.com/ichaly/yugong/core/util"
	"gorm.io/gorm"
	"io"
	"net/http"
	"strconv"
	"time"
)

type DouyinApi struct {
	db     *gorm.DB
	render *base.Render
	spider *douyin.Spider
}

func NewDouyinApi(d *gorm.DB, r *base.Render, s *douyin.Spider) base.Plugin {
	return &DouyinApi{d, r, s}
}

func (my *DouyinApi) Name() string {
	return "DouyinApi"
}

func (my *DouyinApi) Protected() bool {
	return false
}

func (my *DouyinApi) Init(r chi.Router) {
	r.Route("/douyin", func(r chi.Router) {
		r.Get("/start", my.startHandler)
		r.Post("/save", my.saveHandler)
	})
}

func (my *DouyinApi) startHandler(w http.ResponseWriter, r *http.Request) {
	var user data.Douyin
	my.db.First(&user)
	var min int64
	if user.LastTime != nil {
		min = user.LastTime.UnixNano() / 1e6
	}
	min, err := my.spider.GetVideos(user.OpenId, user.Fid, user.Aid, min)
	if err != nil {
		_ = my.render.JSON(w, base.ERROR.WithError(err), base.WithCode(http.StatusBadRequest))
		return
	}
	user.LastTime = util.TimePtr(time.UnixMilli(min))
	//my.db.Save(&user)
	_ = my.render.JSON(w, base.OK.WithData(user))
}

func (my *DouyinApi) saveHandler(w http.ResponseWriter, r *http.Request) {
	var u data.Douyin
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
	info, err := my.spider.GetUserInfo(u.Url)
	if err != nil {
		_ = my.render.JSON(w, base.ERROR.WithError(err))
		return
	}
	u.Fid = info["uid"]
	u.OpenId = info["openid"]
	u.Avatar = info["avatar"]
	u.Nickname = info["nickname"]
	u.From = data.DouYin
	count, err := strconv.ParseInt(info["aweme_count"], 10, 0)
	if err == nil {
		u.ItemCount = count
	}
	my.db.Save(&u)
	_ = my.render.JSON(w, base.OK.WithData(u))
}
