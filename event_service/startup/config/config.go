package config

import (
	"flag"
	cfg "github.com/dislinkt/common/config"
	"os"
)

type Config struct {
	Port        string
	EventDBHost string
	EventDBPort string
	PublicKey   string
}

func NewConfig() *Config {
	devEnv := flag.Bool("dev", false, "use dev environment variables")
	flag.Parse()

	if *devEnv {
		cfg.LoadEnv()
	}

	return &Config{
		Port:        os.Getenv("EVENT_SERVICE_PORT"),
		EventDBHost: os.Getenv("EVENT_DB_HOST"),
		EventDBPort: os.Getenv("EVENT_DB_PORT"),
		PublicKey:   "Dislinkt",
	}
}
