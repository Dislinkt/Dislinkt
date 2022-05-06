package config

import (
	"flag"
	"os"

	cfg "github.com/dislinkt/common/config"
)

type Config struct {
	Port                       string
	AdditionalUserDBHost       string
	AdditionalUserDBPort       string
	NatsHost                   string
	NatsPort                   string
	NatsUser                   string
	NatsPass                   string
	RegisterUserCommandSubject string
	RegisterUserReplySubject   string
}

func NewConfig() *Config {
	devEnv := flag.Bool("dev", false, "use dev environment variables")
	flag.Parse()

	if *devEnv {
		cfg.LoadEnv()
	}

	return &Config{
		Port:                       os.Getenv("ADDITIONAL_USER_SERVICE_PORT"),
		AdditionalUserDBHost:       os.Getenv("ADDITIONAL_USER_DB_HOST"),
		AdditionalUserDBPort:       os.Getenv("ADDITIONAL_USER_DB_PORT"),
		NatsHost:                   os.Getenv("NATS_HOST"),
		NatsPort:                   os.Getenv("NATS_PORT"),
		NatsUser:                   os.Getenv("NATS_USER"),
		NatsPass:                   os.Getenv("NATS_PASS"),
		RegisterUserCommandSubject: os.Getenv("REGISTER_USER_COMMAND_SUBJECT"),
		RegisterUserReplySubject:   os.Getenv("REGISTER_USER_REPLY_SUBJECT"),
	}
}
