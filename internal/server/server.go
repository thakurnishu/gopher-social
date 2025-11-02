package server

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/thakurnishu/gopher-social/internal/config"
	"github.com/thakurnishu/gopher-social/internal/store"
)

type Application struct {
	config config.Config
	store  store.Storage
}

func NewApp(cfg config.Config, store store.Storage) *Application {
	return &Application{
		config: cfg,
		store:  store,
	}
}

func (api *Application) Run(mux *chi.Mux) error {

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
