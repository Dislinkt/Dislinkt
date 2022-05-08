package config

import "os"

type Config struct {
	Port                       string
	Host                       string
	Neo4jUri                   string
	Neo4jUsername              string
	Neo4jPassword              string
	NatsHost                   string
	NatsPort                   string
	NatsUser                   string
	NatsPass                   string
	RegisterUserCommandSubject string
	RegisterUserReplySubject   string
}

func NewConfig() *Config {
	// return &Config{
	//	Port:          "8001",
	//	Host:          "localhost",
	//	Neo4jUri:      "bolt://localhost:7687",
	//	Neo4jUsername: "neo4j",
	//	Neo4jPassword: "password"}

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
	}
}
