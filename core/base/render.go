package base

import (
	"github.com/unrolled/render"
	"net/http"
)

type Option func(code *int)

type Render struct {
	rnd *render.Render
}

func NewRender() *Render {
	return &Render{render.New()}
}

func WithCode(code int) Option {
	return func(c *int) {
		c = &code
	}
}

func (my *Render) JSON(w http.ResponseWriter, v interface{}, opts ...Option) error {
	code := http.StatusOK
	for _, o := range opts {
		o(&code)
	}
	return my.rnd.JSON(w, code, v)
}
