package server

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/thakurnishu/gopher-social/internal/config"
	"github.com/thakurnishu/gopher-social/internal/store"
)

type Server struct {
	config config.Config
	store  store.Storage
}

func NewServer(cfg config.Config, store store.Storage) *Server {
	return &Server{
		config: cfg,
		store:  store,
	}
}

func (api *Server) Run(mux *chi.Mux) error {

	srv := &http.Server{
		Addr:         api.config.Addr,
		Handler:      mux,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}
	log.Printf("server has started at %s", api.config.Addr)
	return srv.ListenAndServe()
}
