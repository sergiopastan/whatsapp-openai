package config

import (
	"os"
	"time"
)

type Config struct {
	DbConfig
}

type DbConfig struct {
	DBHost         string
	DBPoolSize     int
	DBIdlePoolSize int
	DBConnLifetime time.Duration
}

func Load() Config {
	return Config{
		DbConfig: DbConfig{
			DBHost:         os.Getenv("DB_URL"),
			DBPoolSize:     25,
			DBIdlePoolSize: 25,
			DBConnLifetime: 5 * time.Minute,
		},
	}
}
