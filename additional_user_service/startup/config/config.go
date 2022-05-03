package config

import "os"

type Config struct {
	Port                 string
	AdditionalUserDBHost string
	AdditionalUserDBPort string
}

func NewConfig() *Config {
	return &Config{
		Port:                 os.Getenv("ADDITIONAL_USER_SERVICE_PORT"),
		AdditionalUserDBHost: os.Getenv("ADDITIONAL_USER_DB_HOST"),
		AdditionalUserDBPort: os.Getenv("ADDITIONAL_USER_DB_PORT"),
	}
	// return &Config{
	// 	Port:                 "8080",
	// 	AdditionalUserDBHost: "localhost",
	// 	AdditionalUserDBPort: "27017",
	// }
}
