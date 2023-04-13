package base

import (
	"github.com/go-chi/chi/v5"
	"net/http"
)

const (
	apiEndpoint = "/api/"
)

type Server struct {
}

func NewServer() *Server {
	return &Server{}
}

func (my *Server) Attach(r chi.Router) {
	r.Handle(apiEndpoint, my)
}

func (my *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
}
