package config

type Config struct {
	Port               string
	UserHost           string
	UserPort           string
	AdditionalUserHost string
	AdditionalUserPort string
}

func NewConfig() *Config {
	// return &Config{
	// 	Port:     os.Getenv("GATEWAY_PORT"),
	// 	UserHost: os.Getenv("USER_SERVICE_HOST"),
	// 	UserPort: os.Getenv("USER_SERVICE_PORT"),
	// 	AdditionalUserHost: os.Getenv("ADDITIONAL_USER_SERVICE_HOST"),
	// 	AdditionalUserPort: os.Getenv("ADDITIONAL_USER_SERVICE_PORT"),
	// }
	return &Config{
		Port:               "8000",
		UserHost:           "localhost",
		UserPort:           "8090",
		AdditionalUserHost: "localhost",
		AdditionalUserPort: "8090",
	}
}
