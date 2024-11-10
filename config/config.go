package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort string
	DbPath     string
}

func LoadConfig() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("Arquivo .env não encontrado, usando variáveis do ambiente")
	}

	return &Config{
		ServerPort: os.Getenv("SERVER_PORT"),
		DbPath:     os.Getenv("DB_PATH"),
	}
}
