package config

import "os"

var (
	SERVICE_NAME   string
	MONGO_URI      string
	MOGNO_DATABASE string
	JWT_SECRET     string
)

func init() {
	SERVICE_NAME = getEnvOrDefault("SERVICE_NAME", "App")
	MONGO_URI = getEnvOrDefault("MONGO_URI", "mongodb://localhost:27017")
	MOGNO_DATABASE = getEnvOrDefault("MOGNO_DATABASE", "Blog_Default")
	JWT_SECRET = getEnvOrDefault("JWT_SECRET", "RadomString")
}

func getEnvOrDefault(env_key, default_value string) string {

	if v := os.Getenv(env_key); v != "" {
		return v
	}

	return default_value
}
