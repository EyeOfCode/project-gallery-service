package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
    ServerPort     string
    ServerHost     string
    ServerState    string

    MongoDBURI     string
    MongoDBDatabase string
}

func LoadConfig() *Config {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    return &Config{
        ServerPort:      os.Getenv("PORT"),
        ServerHost:      os.Getenv("HOST"),
        ServerState:     os.Getenv("ENV"),

        MongoDBURI:      os.Getenv("MONGO_URI"),
        MongoDBDatabase: os.Getenv("MONGO_DB_NAME"),
    }
}