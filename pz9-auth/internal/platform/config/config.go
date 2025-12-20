package config

import (
	"os"
	"strconv"
)

type Config struct {
	DB_DSN     string
	BcryptCost int
	Addr       string
}

func Load() Config {
	cost := 12

	if v := os.Getenv("BCRYPT_COST"); v != "" {
		integer, err := strconv.Atoi(v)
		if err == nil {
			cost = integer
		}
	}

	addr := "8080"
	if v := os.Getenv("APP_PORT"); v != "" {
		addr = v
	}

	return Config{
		DB_DSN:     os.Getenv("DB_DSN"),
		BcryptCost: cost,
		Addr:       addr,
	}
}
