package serv

import (
	"github.com/bytedance/sonic"
	"github.com/go-chi/chi/v5"
	"github.com/ichaly/jingwei/core/base"
	"github.com/ichaly/jingwei/core/data"
	"gorm.io/gorm"
	"io"
	"net/http"
)

type UserService struct {
	db     *gorm.DB
	render *base.Render
}

func NewUserService(d *gorm.DB, r *base.Render) base.Plugin {
	return &UserService{d, r}
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
	my.db.Save(&u)
	_ = my.render.JSON(w, base.OK.WithData(u))
}
