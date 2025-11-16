package main

import (
	"log"

	"github.com/thakurnishu/gopher-social/internal/config"
	"github.com/thakurnishu/gopher-social/internal/db"
	"github.com/thakurnishu/gopher-social/internal/store"
)


func main() {
	cfg := config.Load()

	dbConn, err := db.New(
		cfg.DB.Addr,
		cfg.DB.MaxOpenConns,
		cfg.DB.MaxIdleConns,
		cfg.DB.MaxIdleTime,
	)
	if err != nil {
		log.Panic(err)
	}
	defer dbConn.Close()
	log.Println("database connection establised")
	storage := store.NewStorage(dbConn)
	
	db.Seed(storage)
}
