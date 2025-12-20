package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port          string
	PublicRsaKey  string
	PrivateRsaKey string
	AccessTTL     time.Duration
	RefreshTTL    time.Duration
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int
}

func Load() Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	PublicRsaKey := os.Getenv("PUBLIC_RSA_KEY")
	if PublicRsaKey == "" {
		log.Fatal("PUBLIC_RSA_KEY is required")
	}

	PrivateRsaKey := os.Getenv("PRIVATE_RSA_KEY")
	if PrivateRsaKey == "" {
		log.Fatal("PRIVATE_RSA_KEY is required")
	}

	accessTTL := os.Getenv("ACCESS_TTL")
	if accessTTL == "" {
		accessTTL = "15m"
	}
	accessTTLDur, err := time.ParseDuration(accessTTL)
	if err != nil {
		log.Fatal("bad ACCESS_TTL")
	}

	refreshTTL := os.Getenv("REFRESH_TTL")
	if refreshTTL == "" {
		refreshTTL = "168h"
	}
	refreshTTLDur, err := time.ParseDuration(refreshTTL)
	if err != nil {
		log.Fatal("bad REFRESH_TTL")
	}

	redisHost := os.Getenv("REDIS_HOST")
	if redisHost == "" {
		redisHost = "localhost"
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")

	redisPort := os.Getenv("REDIS_PORT")
	if redisPort == "" {
		redisPort = "6379"
	}

	redisDBStr := os.Getenv("REDIS_DB")
	if redisDBStr == "" {
		redisDBStr = "0"
	}
	redisDB, err := strconv.Atoi(redisDBStr)
	if err != nil {
		log.Printf("Invalid REDIS_DB value '%s', using default 0", redisDBStr)
		redisDB = 0
	}

	return Config{
		Port:          ":" + port,
		PublicRsaKey:  PublicRsaKey,
		PrivateRsaKey: PrivateRsaKey,
		AccessTTL:     accessTTLDur,
		RefreshTTL:    refreshTTLDur,
		RedisHost:     redisHost,
		RedisPort:     redisPort,
		RedisPassword: redisPassword,
		RedisDB:       redisDB,
	}
}
