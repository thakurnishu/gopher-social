package main

import (
	"log"

	"github.com/thakurnishu/gopher-social/internal/config"
	"github.com/thakurnishu/gopher-social/internal/db"
	"github.com/thakurnishu/gopher-social/internal/server"
	"github.com/thakurnishu/gopher-social/internal/store"
)

func main() {
	cfg := config.Load()

	db, err := db.New(
		cfg.DB.Addr,
		cfg.DB.MaxOpenConns,
		cfg.DB.MaxIdleConns,
		cfg.DB.MaxIdleTime,
	)
	if err != nil {
		log.Panic(err)
	}
	defer db.Close()
	log.Println("database connection establised")
	storage := store.NewStorage(db)

	app := server.NewServer(cfg, storage)
	mux := app.RegisterRoutes()
	log.Fatal(app.Run(mux))
}
