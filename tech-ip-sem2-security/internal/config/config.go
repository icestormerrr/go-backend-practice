package config

import "os"

type Config struct {
	Addr     string
	CertFile string
	KeyFile  string
	DSN      string
}

func New() Config {
	return Config{
		Addr:     getEnv("HTTPS_ADDR", ":8443"),
		CertFile: getEnv("TLS_CERT_FILE", "certs/server.crt"),
		KeyFile:  getEnv("TLS_KEY_FILE", "certs/server.key"),
		DSN:      getEnv("DATABASE_DSN", "postgres://postgres:postgres@localhost:5432/study_security?sslmode=disable"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
