package config

import (
	"flag"
	cfg "github.com/dislinkt/common/config"
	"os"
)

type Config struct {
	Port               string
	NotificationDBHost string
	NotificationDBPort string
	PublicKey          string
}

func NewConfig() *Config {
	devEnv := flag.Bool("dev", false, "use dev environment variables")
	flag.Parse()

	if *devEnv {
		cfg.LoadEnv()
	}

	return &Config{
		Port:               os.Getenv("NOTIFICATION_SERVICE_PORT"),
		NotificationDBHost: os.Getenv("NOTIFICATION_DB_HOST"),
		NotificationDBPort: os.Getenv("NOTIFICATION_DB_PORT"),
		PublicKey:          "Dislinkt",
	}
}
