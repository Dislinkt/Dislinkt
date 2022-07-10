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
	AddEducationCommandSubject string
	AddEducationReplySubject   string
	AddSkillCommandSubject     string
	AddSkillReplySubject       string
	DeleteSkillCommandSubject  string
	DeleteSkillReplySubject    string
	UpdateSkillCommandSubject  string
	UpdateSkillReplySubject    string
	PublicKey                  string
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
		AddEducationCommandSubject: os.Getenv("ADD_EDUCATION_COMMAND_SUBJECT"),
		AddEducationReplySubject:   os.Getenv("ADD_EDUCATION_REPLY_SUBJECT"),
		AddSkillCommandSubject:     os.Getenv("ADD_SKILL_COMMAND_SUBJECT"),
		AddSkillReplySubject:       os.Getenv("ADD_SKILL_REPLY_SUBJECT"),
		DeleteSkillCommandSubject:  os.Getenv("DELETE_SKILL_COMMAND_SUBJECT"),
		DeleteSkillReplySubject:    os.Getenv("DELETE_SKILL_REPLY_SUBJECT"),
		UpdateSkillCommandSubject:  os.Getenv("UPDATE_SKILL_COMMAND_SUBJECT"),
		UpdateSkillReplySubject:    os.Getenv("UPDATE_SKILL_REPLY_SUBJECT"),
		PublicKey:                  "Dislinkt",
	}
}
