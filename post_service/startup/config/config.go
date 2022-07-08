package config

import (
	"flag"
	"os"

	cfg "github.com/dislinkt/common/config"
)

type Config struct {
	Port                         string
	PostDBHost                   string
	PostDBPort                   string
	NatsHost                     string
	NatsPort                     string
	NatsUser                     string
	NatsPass                     string
	RegisterUserCommandSubject   string
	RegisterUserReplySubject     string
	UpdateUserCommandSubject     string
	UpdateUserReplySubject       string
	CreateJobOfferCommandSubject string
	CreateJobOfferReplySubject   string
	PublicKey                    string
}

func NewConfig() *Config {
	devEnv := flag.Bool("dev", false, "use dev environment variables")
	flag.Parse()

	if *devEnv {
		cfg.LoadEnv()
	}

	return &Config{
		Port:                         os.Getenv("POST_SERVICE_PORT"),
		PostDBHost:                   os.Getenv("POST_DB_HOST"),
		PostDBPort:                   os.Getenv("POST_DB_PORT"),
		NatsHost:                     os.Getenv("NATS_HOST"),
		NatsPort:                     os.Getenv("NATS_PORT"),
		NatsUser:                     os.Getenv("NATS_USER"),
		NatsPass:                     os.Getenv("NATS_PASS"),
		RegisterUserCommandSubject:   os.Getenv("REGISTER_USER_COMMAND_SUBJECT"),
		RegisterUserReplySubject:     os.Getenv("REGISTER_USER_REPLY_SUBJECT"),
		UpdateUserCommandSubject:     os.Getenv("UPDATE_USER_COMMAND_SUBJECT"),
		UpdateUserReplySubject:       os.Getenv("UPDATE_USER_REPLY_SUBJECT"),
		CreateJobOfferCommandSubject: os.Getenv("CREATE_JOB_COMMAND_SUBJECT"),
		CreateJobOfferReplySubject:   os.Getenv("CREATE_JOB_REPLY_SUBJECT"),
		PublicKey:                    "Dislinkt",
	}
}
