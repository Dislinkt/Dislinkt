package config

import (
	"flag"
	cfg "github.com/dislinkt/common/config"
)

type Config struct {
	Port               string
	UserHost           string
	UserPort           string
	AdditionalUserHost string
	AdditionalUserPort string
	ConnectionHost     string
	ConnectionPort     string
}

func NewConfig() *Config {
	devEnv := flag.Bool("dev", false, "use dev environment variables")
	flag.Parse()

	if *devEnv {
		cfg.LoadEnv()
	}
	return &Config{
		Port:           "8080",
		ConnectionHost: "localhost",
		ConnectionPort: "8001",
		//UserHost:           os.Getenv("USER_SERVICE_HOST"),
		//UserPort:           os.Getenv("USER_SERVICE_PORT"),
		//AdditionalUserHost: os.Getenv("ADDITIONAL_USER_SERVICE_HOST"),
		//AdditionalUserPort: os.Getenv("ADDITIONAL_USER_SERVICE_PORT"),
	}

}
