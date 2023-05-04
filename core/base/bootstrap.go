package base

import (
	"context"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/fx"
	"net/http"
)

func Bootstrap(
	l fx.Lifecycle, c *Config, h *chi.Mux, g PluginGroup, s *Server,
) {
	//init middlewares
	for _, m := range g.Middlewares {
		if !m.Protected() {
			m.Init(h)
		}
	}
	//wrap service
	h.Group(func(h chi.Router) {
		for _, m := range g.Middlewares {
			if m.Protected() {
				m.Init(h)
			}
		}
		s.Attach(h)
		for _, p := range g.Plugins {
			if p.Protected() {
				p.Init(h)
			}
		}
	})
	//init plugins
	for _, p := range g.Plugins {
		if !p.Protected() {
			p.Init(h)
		}
	}
	srv := &http.Server{Addr: fmt.Sprintf(":%v", c.App.Port), Handler: h}
	l.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				err := srv.ListenAndServe()
				if err != nil {
					fmt.Printf("%v failed to start: %v", c.App.Name, err)
					return
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			fmt.Printf("%v shutdown complete", c.App.Name)
			return srv.Shutdown(ctx)
		},
	})
}
