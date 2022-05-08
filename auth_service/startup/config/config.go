package config

import "os"

type Config struct {
	Port       string
	UserDBHost string
	UserDBPort string
	UserDBName string
	UserDBUser string
	UserDBPass string
	Secret     string
}

func NewConfig() *Config {
	return &Config{
		Port:       os.Getenv("AUTH_SERVICE_PORT"),
		UserDBHost: os.Getenv("AUTH_DB_HOST"),
		UserDBPort: os.Getenv("AUTH_DB_PORT"),
		UserDBName: os.Getenv("AUTH_DB_NAME"),
		UserDBUser: os.Getenv("AUTH_DB_USER"),
		UserDBPass: os.Getenv("AUTH_DB_PASS"),
		Secret:     os.Getenv("SECRET"),
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
