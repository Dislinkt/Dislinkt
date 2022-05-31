package config

import (
	"os"
)

type Config struct {
	Port                       string
	AuthDBHost                 string
	AuthDBPort                 string
	AuthDBName                 string
	AuthDBUser                 string
	AuthDBPass                 string
	Secret                     string
	NatsHost                   string
	NatsPort                   string
	NatsUser                   string
	NatsPass                   string
	RegisterUserCommandSubject string
	RegisterUserReplySubject   string
	EmailSender                string
	EmailPassword              string
	EmailHost                  string
	EmailPort                  string
	PublicKey                  string
}

func NewConfig() *Config {
	// devEnv := flag.Bool("dev", false, "use dev environment variables")
	// flag.Parse()
	//
	// if *devEnv {
	// 	cfg.LoadEnv()
	// }

	return &Config{
		Port:                       os.Getenv("AUTH_SERVICE_PORT"),
		AuthDBHost:                 os.Getenv("AUTH_DB_HOST"),
		AuthDBPort:                 os.Getenv("AUTH_DB_PORT"),
		AuthDBName:                 os.Getenv("AUTH_DB_NAME"),
		AuthDBUser:                 os.Getenv("AUTH_DB_USER"),
		AuthDBPass:                 os.Getenv("AUTH_DB_PASS"),
		Secret:                     os.Getenv("SECRET"),
		NatsHost:                   os.Getenv("NATS_HOST"),
		NatsPort:                   os.Getenv("NATS_PORT"),
		NatsUser:                   os.Getenv("NATS_USER"),
		NatsPass:                   os.Getenv("NATS_PASS"),
		RegisterUserCommandSubject: os.Getenv("REGISTER_USER_COMMAND_SUBJECT"),
		RegisterUserReplySubject:   os.Getenv("REGISTER_USER_REPLY_SUBJECT"),
		EmailSender:                "xmlbsep@gmail.com",
		EmailHost:                  "smtp.gmail.com",
		EmailPassword:              "sifra123",
		EmailPort:                  "587",
		PublicKey:                  "Dislinkt",
	}
	// return &Config{
	// 	Port:       "8090",
	// 	UserDBHost: "localhost",
	// 	UserDBPort: "5432",
	// 	UserDBName: "dislinkt-auth",
	// 	UserDBUser: "postgres",
	// 	UserDBPass: "password",
	// }
}
