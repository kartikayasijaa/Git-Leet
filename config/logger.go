package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func InitLogger() *log.Logger {
	logger := log.New(os.Stdout, "Log: ", log.Lshortfile)

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading environment variables")
	}

	return logger
}
