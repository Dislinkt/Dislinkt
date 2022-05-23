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
	EmailSender                string
	EmailPassword              string
	EmailHost                  string
	EmailPort                  string
	PublicKey                  string
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
		EmailSender:                os.Getenv("EMAIL_SENDER"),
		EmailHost:                  os.Getenv("EMAIL_HOST"),
		EmailPassword:              os.Getenv("EMAIL_PASSWORD"),
		EmailPort:                  os.Getenv("EMAIL_PORT"),
		PublicKey:                  "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA0AzWYJTc9jiPn+RMNjMJ\nhscn8hg/Mt0U22efM6IvM83CyQCiFHP1Z8rs2HFqRbid/hQxW23HrXQzKx5hGPdU\n14ncF8oN7utDQxdq6ivTsF1tMQtHWb2jnYmpKwTyelbMMGKLHj3yy2j59Y/X94EX\nPNtQtgAO9FF5gKzjkaBu6KzLU2RJC9bADVd5sotM/JP/Ce5D/97XV7i1KStTUDiV\nfDBWCkDylBTQTmI1rO9MdayVduuAzNdWXRfyqKcWI2i4pA1aaskiaViVsIhF3ksm\nYW4Bu0RxK5SP2byHj7pv93XsabA+QXZ37QRhYzBxx6nS0x/dNtAxIltIBZaeSTN0\ngQIDAQAB\n-----END PUBLIC KEY-----",
	}
}
