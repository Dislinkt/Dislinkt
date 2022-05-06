package config

import (
	"flag"
	"os"

	cfg "github.com/dislinkt/common/config"
)

type Config struct {
	Port                       string
	UserDBHost                 string
	UserDBPort                 string
	UserDBName                 string
	UserDBUser                 string
	UserDBPass                 string
	NatsHost                   string
	NatsPort                   string
	NatsUser                   string
	NatsPass                   string
	RegisterUserCommandSubject string
	RegisterUserReplySubject   string
	JaegerServiceName          string
}

func NewConfig() *Config {
	devEnv := flag.Bool("dev", false, "use dev environment variables")
	flag.Parse()

	if *devEnv {
		cfg.LoadEnv()
	}

	return &Config{
		Port:                       os.Getenv("USER_SERVICE_PORT"),
		UserDBHost:                 os.Getenv("USER_DB_HOST"),
		UserDBPort:                 os.Getenv("USER_DB_PORT"),
		UserDBName:                 os.Getenv("USER_DB_NAME"),
		UserDBUser:                 os.Getenv("USER_DB_USER"),
		UserDBPass:                 os.Getenv("USER_DB_PASS"),
		NatsHost:                   os.Getenv("NATS_HOST"),
		NatsPort:                   os.Getenv("NATS_PORT"),
		NatsUser:                   os.Getenv("NATS_USER"),
		NatsPass:                   os.Getenv("NATS_PASS"),
		RegisterUserCommandSubject: os.Getenv("REGISTER_USER_COMMAND_SUBJECT"),
		RegisterUserReplySubject:   os.Getenv("REGISTER_USER_REPLY_SUBJECT"),
		JaegerServiceName:          os.Getenv("JAEGER_SERVICE_NAME"),
	}
}
