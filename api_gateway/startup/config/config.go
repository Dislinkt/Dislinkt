package config

import (
	"flag"
	"os"

	cfg "github.com/dislinkt/common/config"
)

type Config struct {
	Port               string
	UserHost           string
	UserPort           string
	AuthHost           string
	AuthPort           string
	AdditionalUserHost string
	AdditionalUserPort string
	ConnectionHost     string
	ConnectionPort     string
	PostHost           string
	PostPort           string
}

func NewConfig() *Config {
	devEnv := flag.Bool("dev", false, "use dev environment variables")
	flag.Parse()

	if *devEnv {
		cfg.LoadEnv()
	}
	return &Config{
		Port:               os.Getenv("GATEWAY_PORT"),
		UserHost:           os.Getenv("USER_SERVICE_HOST"),
		UserPort:           os.Getenv("USER_SERVICE_PORT"),
		AuthHost:           os.Getenv("AUTH_SERVICE_HOST"),
		AuthPort:           os.Getenv("AUTH_SERVICE_PORT"),
		ConnectionHost:     os.Getenv("CONNECTION_SERVICE_HOST"),
		ConnectionPort:     os.Getenv("CONNECTION_SERVICE_PORT"),
		PostHost:           os.Getenv("POST_SERVICE_HOST"),
		PostPort:           os.Getenv("POST_SERVICE_PORT"),
		AdditionalUserHost: os.Getenv("ADDITIONAL_USER_SERVICE_HOST"),
		AdditionalUserPort: os.Getenv("ADDITIONAL_USER_SERVICE_PORT"),
	}

}
