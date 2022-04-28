package config

type Config struct {
	Port              string
	UserDBHost        string
	UserDBPort        string
	UserDBName        string
	UserDBUser        string
	UserDBPass        string
	JaegerServiceName string
}

func NewConfig() *Config {
	// return &Config{
	// 	Port:              os.Getenv("USER_SERVICE_PORT"),
	// 	UserDBHost:        os.Getenv("USER_DB_HOST"),
	// 	UserDBPort:        os.Getenv("USER_DB_PORT"),
	// 	UserDBName:        os.Getenv("USER_DB_NAME"),
	// 	UserDBUser:        os.Getenv("USER_DB_USER"),
	// 	UserDBPass:        os.Getenv("USER_DB_PASS"),
	// 	JaegerServiceName: os.Getenv("JAEGER_SERVICE_NAME"),
	// }
	return &Config{
		Port:       "8090",
		UserDBHost: "localhost",
		UserDBPort: "5432",
		UserDBName: "dislinkt-users",
		UserDBUser: "postgres",
		UserDBPass: "password",
	}
}
