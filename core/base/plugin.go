package base

import (
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
)

type PluginGroup struct {
	fx.In
	Plugins     []Plugin `group:"plugin"`
	Middlewares []Plugin `group:"middleware"`
}

type Plugin interface {
	Name() string
	Protected() bool
	Init(r chi.Router)
}
