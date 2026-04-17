package config

import "os"

type Config struct {
	DBUrl string
}

func Load() Config {
	return Config{
		DBUrl: os.Getenv("DB_URL"),
	}
}