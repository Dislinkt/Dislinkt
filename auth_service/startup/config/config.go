package config

import (
	"flag"
	"os"

	cfg "github.com/dislinkt/common/config"
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
	PublicKey                  string
}

func NewConfig() *Config {
	devEnv := flag.Bool("dev", false, "use dev environment variables")
	flag.Parse()

	if *devEnv {
		cfg.LoadEnv()
	}

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
		PublicKey:                  "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0AzWYJTc9jiPn+RMNjMJ\nhscn8hg/Mt0U22efM6IvM83CyQCiFHP1Z8rs2HFqRbid/hQxW23HrXQzKx5hGPdU\n14ncF8oN7utDQxdq6ivTsF1tMQtHWb2jnYmpKwTyelbMMGKLHj3yy2j59Y/X94EX\nPNtQtgAO9FF5gKzjkaBu6KzLU2RJC9bADVd5sotM/JP/Ce5D/97XV7i1KStTUDiV\nfDBWCkDylBTQTmI1rO9MdayVduuAzNdWXRfyqKcWI2i4pA1aaskiaViVsIhF3ksm\nYW4Bu0RxK5SP2byHj7pv93XsabA+QXZ37QRhYzBxx6nS0x/dNtAxIltIBZaeSTN0\ngQIDAQAB\n-----END PUBLIC KEY-----",
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
