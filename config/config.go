package config

type Config struct {
	ServiceName      string
	ServiceHost      string
	ServicePort      string
	PostgresHost     string
	PostgresPort     int
	PostgresUser     string
	PostgresPassword string
	PostgresDatabase string
	DefaultOffset    string
	DefaultLimit     string
}

func Load() Config {
	config := Config{}

	config.ServiceName = "img-service"
	config.ServiceHost = "localhost"
	config.ServicePort = ":9000"
	config.PostgresHost = "localhost"
	config.PostgresPort = 5432
	config.PostgresUser = "postgres"
	config.PostgresPassword = "postgres"
	config.PostgresDatabase = "img_service"
	config.DefaultOffset = "0"
	config.DefaultLimit = "10"

	return config
}
