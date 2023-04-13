package plug

import (
	"github.com/bytedance/sonic"
	"github.com/go-chi/chi/v5"
	"github.com/ichaly/jingwei/core/base"
	"github.com/ichaly/jingwei/core/data"
	"github.com/ichaly/jingwei/core/serv"
	"gorm.io/gorm"
	"io"
	"net/http"
	"strconv"
)

type UserService struct {
	db     *gorm.DB
	render *base.Render
	spider *serv.Spider
}

func NewUserService(d *gorm.DB, r *base.Render, s *serv.Spider) base.Plugin {
	return &UserService{d, r, s}
}

func (my *UserService) Name() string {
	return "UserService"
}

func (my *UserService) Protected() bool {
	return false
}

func (my *UserService) Init(r chi.Router) {
	r.Route("/user", func(r chi.Router) {
		r.Post("/save", my.ServeHTTP)
	})
}

func (my *UserService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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
	u.Did = info["did"]
	u.Nickname = info["nickname"]
	u.Avatar = info["avatar"]
	count, err := strconv.ParseInt(info["aweme_count"], 10, 0)
	if err != nil {
		u.ItemCount = count
	}
	my.db.Save(&u)
	_ = my.render.JSON(w, base.OK.WithData(u))
}
