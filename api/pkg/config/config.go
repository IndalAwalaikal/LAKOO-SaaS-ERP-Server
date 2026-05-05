package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port           string
	DBHost         string
	DBUser         string
	DBPassword     string
	DBName         string
	RedisHost      string
	RedisPort      string
	JWTSecret      string
	MinioEndpoint  string
	MinioAccessKey string
	MinioSecretKey string
	MinioBucket    string
	MinioUseSSL    bool
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, relying on environment variables")
	}

	return &Config{
		Port:           getEnv("PORT", "8080"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBUser:         getEnv("DB_USER", "lakoo_user"),
		DBPassword:     getEnv("DB_PASSWORD", "lakoo_secret"),
		DBName:         getEnv("DB_NAME", "lakoo_db"),
		RedisHost:      getEnv("REDIS_HOST", "localhost"),
		RedisPort:      getEnv("REDIS_PORT", "6379"),
		JWTSecret:      getEnv("JWT_SECRET", "super_secret_jwt_key_for_lakoo_saas"),
		MinioEndpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinioAccessKey: getEnv("MINIO_ROOT_USER", "minioadmin"),
		MinioSecretKey: getEnv("MINIO_ROOT_PASSWORD", "minioadmin"),
		MinioBucket:    getEnv("MINIO_BUCKET", "lakoo-media"),
		MinioUseSSL:    getEnv("MINIO_USE_SSL", "false") == "true",
	}
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
