package env

import (
	"os"
	"strconv"
	"time"
)

func GetString(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}
	// log.Println("got string from env:", key)
	return val
}

func GetInt(key string, fallback int) int {
	val, ok := os.LookupEnv(key)
	if !ok {
		return fallback
	}

	valAsInt, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	// log.Println("got int from env:", key)
	return valAsInt
}

func GetDuration(key, fallback string) time.Duration {
	val, ok := os.LookupEnv(key)
	if !ok {
		fallbackDuration, err := time.ParseDuration(fallback)
		if err != nil {
			return time.Minute * 15
		}
		return fallbackDuration
	}

	duration, err := time.ParseDuration(val)
	if err != nil {
		return time.Minute * 15
	}
	
	// log.Println("got duration from env:", key)
	return duration
}
