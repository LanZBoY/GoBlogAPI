package config

import "os"

var (
	MONGO_URI string
)

func init() {
	MONGO_URI = getEnvOrDefault("MONGO_URI", "mongodb://localhost:27017")
}

func getEnvOrDefault(env_key, default_value string) string {

	if v := os.Getenv(env_key); v != "" {
		return v
	}

	return default_value
}
