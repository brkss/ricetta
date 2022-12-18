package utils

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBDriver      string
	DBSource      string
	ServerAddress string
}

func LoadConfig() (*Config, error) {

	err := godotenv.Load()
	if err != nil {
		return nil, err
	}
	config := &Config{
		DBDriver:      os.Getenv("DB_DRIVER"),
		DBSource:      os.Getenv("DB_SOURCE"),
		ServerAddress: os.Getenv("SERVER_ADDRESS"),
	}
	return config, nil
}
