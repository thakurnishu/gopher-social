package config

import (
	"log"
	"time"

	"github.com/joho/godotenv"
	"github.com/thakurnishu/gopher-social/internal/env"
)

type Config struct {
	Addr    string
	DB      DBConfig
	Env     string
	Version string
}

type DBConfig struct {
	Addr         string
	MaxOpenConns int
	MaxIdleConns int
	MaxIdleTime  time.Duration
}

func Load() Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := Config{
		Addr: env.GetString("SERVER_PORT", ":8080"),
		DB: DBConfig{
			Addr:         env.GetString("DB_ADDR", "postgres://admin:adminpassword@localhost/gopher_social?sslmode=disable"),
			MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
			MaxIdleTime:  env.GetDuration("DB_MAX_IDLE_TIME", "15m"),
		},
		Env:     env.GetString("ENVIRNMENT", "development"),
		Version: env.GetString("VERSION", "0.0.1"),
	}

	return cfg
}
