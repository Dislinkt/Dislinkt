package config

import (
	"flag"
	"os"

	cfg "github.com/dislinkt/common/config"
)

type Config struct {
	Port              string
	UserDBHost        string
	UserDBPort        string
	UserDBName        string
	UserDBUser        string
	UserDBPass        string
	JaegerServiceName string
}

func NewConfig() *Config {
	devEnv := flag.Bool("dev", false, "use dev environment variables")
	flag.Parse()

	if *devEnv {
		cfg.LoadEnv()
	}

	return &Config{
		Port:              os.Getenv("USER_SERVICE_PORT"),
		UserDBHost:        os.Getenv("USER_DB_HOST"),
		UserDBPort:        os.Getenv("USER_DB_PORT"),
		UserDBName:        os.Getenv("USER_DB_NAME"),
		UserDBUser:        os.Getenv("USER_DB_USER"),
		UserDBPass:        os.Getenv("USER_DB_PASS"),
		JaegerServiceName: os.Getenv("JAEGER_SERVICE_NAME"),
	}
}
