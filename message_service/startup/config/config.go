package config

import (
	"flag"
	cfg "github.com/dislinkt/common/config"
	"os"
)

type Config struct {
	Port          string
	MessageDBHost string
	MessageDBPort string
	PublicKey     string
}

func NewConfig() *Config {
	devEnv := flag.Bool("dev", false, "use dev environment variables")
	flag.Parse()

	if *devEnv {
		cfg.LoadEnv()
	}

	return &Config{
		Port:          os.Getenv("MESSAGE_SERVICE_PORT"),
		MessageDBHost: os.Getenv("MESSAGE_DB_HOST"),
		MessageDBPort: os.Getenv("MESSAGE_DB_PORT"),
		PublicKey:     "Dislinkt",
	}
}
