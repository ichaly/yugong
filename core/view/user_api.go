package view

import (
	"github.com/bytedance/sonic"
	"github.com/go-chi/chi/v5"
	"github.com/ichaly/yugong/core/base"
	"github.com/ichaly/yugong/core/data"
	"github.com/ichaly/yugong/core/serv"
	"github.com/ichaly/yugong/core/util"
	"gorm.io/gorm"
	"io"
	"net/http"
	"strconv"
	"time"
)

type UserApi struct {
	db     *gorm.DB
	render *base.Render
	spider *serv.Spider
}

func NewUserApi(d *gorm.DB, r *base.Render, s *serv.Spider) base.Plugin {
	return &UserApi{d, r, s}
}

func (my *UserApi) Name() string {
	return "UserApi"
}

func (my *UserApi) Protected() bool {
	return false
}

func (my *UserApi) Init(r chi.Router) {
	r.Route("/user", func(r chi.Router) {
		r.Get("/start", my.startHandler)
		r.Post("/save", my.saveHandler)
	})
}

func (my *UserApi) startHandler(w http.ResponseWriter, r *http.Request) {
	var user data.User
	my.db.First(&user)
	var min int64
	if user.LastVisit != nil {
		min = user.LastVisit.UnixNano() / 1e6
	}
	min, err := my.spider.GetVideos(user.OpenId, user.Did, user.Aid, min)
	if err != nil {
		_ = my.render.JSON(w, base.ERROR.WithError(err), base.WithCode(http.StatusBadRequest))
		return
	}
	user.LastVisit = util.TimePtr(time.UnixMilli(min))
	//my.db.Save(&user)
	_ = my.render.JSON(w, base.OK.WithData(user))
}

func (my *UserApi) saveHandler(w http.ResponseWriter, r *http.Request) {
	var u data.User
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
	u.Did = info["uid"]
	u.OpenId = info["openid"]
	u.Avatar = info["avatar"]
	u.Nickname = info["nickname"]
	count, err := strconv.ParseInt(info["aweme_count"], 10, 0)
	if err == nil {
		u.ItemCount = count
	}
	my.db.Save(&u)
	_ = my.render.JSON(w, base.OK.WithData(u))
}
