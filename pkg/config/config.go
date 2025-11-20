package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port  string
	DBURL string
}

func Load() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	port := os.Getenv("APP_PORT")
	if port == "" {
		port = "8080"
	}

	dbURL := os.Getenv("DB_CONN")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is required")
	}

	return &Config{
		Port:  port,
		DBURL: dbURL,
	}
}
