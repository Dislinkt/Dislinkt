package config

import (
	"flag"
	cfg "github.com/dislinkt/common/config"
	"os"
)

type Config struct {
	Port       string
	PostDBHost string
	PostDBPort string
}

func NewConfig() *Config {
	devEnv := flag.Bool("dev", false, "use dev environment variables")
	flag.Parse()

	if *devEnv {
		cfg.LoadEnv()
	}

	return &Config{
		Port:       os.Getenv("POST_SERVICE_PORT"),
		PostDBHost: os.Getenv("POST_DB_HOST"),
		PostDBPort: os.Getenv("POST_DB_PORT"),
	}
}
