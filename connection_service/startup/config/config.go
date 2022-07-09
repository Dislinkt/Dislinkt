package config

import (
	"flag"
	cfg "github.com/dislinkt/common/config"
	"os"
)

type Config struct {
	Port                         string
	Host                         string
	Neo4jUri                     string
	Neo4jUsername                string
	Neo4jPassword                string
	NatsHost                     string
	NatsPort                     string
	NatsUser                     string
	NatsPass                     string
	RegisterUserCommandSubject   string
	RegisterUserReplySubject     string
	PatchUserCommandSubject      string
	PatchUserReplySubject        string
	CreateJobOfferCommandSubject string
	CreateJobOfferReplySubject   string
	AddEducationCommandSubject   string
	AddEducationReplySubject     string
	AddSkillCommandSubject       string
	AddSkillReplySubject         string
	PublicKey                    string
}

func NewConfig() *Config {
	// return &Config{
	//	Port:          "8001",
	//	Host:          "localhost",
	//	Neo4jUri:      "bolt://localhost:7687",
	//	Neo4jUsername: "neo4j",
	//	Neo4jPassword: "password"}

	devEnv := flag.Bool("dev", false, "use dev environment variables")
	flag.Parse()

	if *devEnv {
		cfg.LoadEnv()
	}

	return &Config{
		Port:                       os.Getenv("CONNECTION_SERVICE_PORT"),
		Neo4jUri:                   "neo4j://neo4j:7687",
		Neo4jUsername:              "neo4j",
		Neo4jPassword:              "password",
		NatsHost:                   os.Getenv("NATS_HOST"),
		NatsPort:                   os.Getenv("NATS_PORT"),
		NatsUser:                   os.Getenv("NATS_USER"),
		NatsPass:                   os.Getenv("NATS_PASS"),
		RegisterUserCommandSubject: os.Getenv("REGISTER_USER_COMMAND_SUBJECT"),
		RegisterUserReplySubject:   os.Getenv("REGISTER_USER_REPLY_SUBJECT"),
		PatchUserCommandSubject:    os.Getenv("PATCH_USER_COMMAND_SUBJECT"),
		PatchUserReplySubject:      os.Getenv("PATCH_USER_REPLY_SUBJECT"),
		//CreateJobOfferCommandSubject: "job.create.command",
		//CreateJobOfferReplySubject:   "job.create.reply",
		CreateJobOfferCommandSubject: os.Getenv("CREATE_JOB_COMMAND_SUBJECT"),
		CreateJobOfferReplySubject:   os.Getenv("CREATE_JOB_REPLY_SUBJECT"),
		AddEducationCommandSubject:   os.Getenv("ADD_EDUCATION_COMMAND_SUBJECT"),
		AddEducationReplySubject:     os.Getenv("ADD_EDUCATION_REPLY_SUBJECT"),
		AddSkillCommandSubject:       os.Getenv("ADD_SKILL_COMMAND_SUBJECT"),
		AddSkillReplySubject:         os.Getenv("ADD_SKILL_REPLY_SUBJECT"),
		PublicKey:                    "Dislinkt",
	}
}
