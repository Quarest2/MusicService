package config

import (
	"os"
	"strconv"
)

var AppConfig *Config

func LoadConfig() {
	AppConfig = &Config{
		Port: getEnv("PORT", "8080"),

		MinioEndpoint:     getEnv("MINIO_ENDPOINT", "localhost:9000"),
		BucketName:        getEnv("MINIO_BUCKET_NAME", "defaultBucket"),
		MinioRootUser:     getEnv("MINIO_ROOT_USER", "root"),
		MinioRootPassword: getEnv("MINIO_ROOT_PASSWORD", "minio_password"),
		MinioUseSSL:       getEnvAsBool("MINIO_USE_SSL", false),
	}
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if valueStr := getEnv(key, ""); valueStr != "" {
		if value, err := strconv.ParseBool(valueStr); err == nil {
			return value
		}
	}
	return defaultValue
}
