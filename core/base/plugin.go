package base

import "github.com/go-chi/chi/v5"

type Plugin interface {
	Name() string
	Protected() bool
	Init(r chi.Router)
}
