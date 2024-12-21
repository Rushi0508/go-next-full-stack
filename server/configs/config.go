package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	MongoURI    string
	DBName      string
	JWTSecret   string
	Environment string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	return &Config{
		Port:        getEnvOrDefault("PORT", "8080"),
		MongoURI:    getEnvOrDefault("MONGO_URI", "mongodb://localhost:27017"),
		DBName:      getEnvOrDefault("DB_NAME", "blog_db"),
		JWTSecret:   getEnvOrDefault("JWT_SECRET", "your-secret-key"),
		Environment: getEnvOrDefault("ENVIRONMENT", "development"),
	}
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
