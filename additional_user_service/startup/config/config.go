package config

import (
	"flag"
	"os"

	cfg "github.com/dislinkt/common/config"
)

type Config struct {
	Port                 string
	AdditionalUserDBHost string
	AdditionalUserDBPort string
}

func NewConfig() *Config {
	devEnv := flag.Bool("dev", false, "use dev environment variables")
	flag.Parse()

	if *devEnv {
		cfg.LoadEnv()
	}

	return &Config{
		Port:                 os.Getenv("ADDITIONAL_USER_SERVICE_PORT"),
		AdditionalUserDBHost: os.Getenv("ADDITIONAL_USER_DB_HOST"),
		AdditionalUserDBPort: os.Getenv("ADDITIONAL_USER_DB_PORT"),
	}
}
