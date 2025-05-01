package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type Config struct {
	DBHost           string
	DBPort           string
	DBUser           string
	DBName           string
	DBPassword       string
	RedisAddr        string
	RedisPassword    string
	RedisDB          int
	TelegramBotToken string
	Port             string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		return nil, err
	}

	return &Config{
		DBHost:           os.Getenv("DB_HOST"),
		DBPort:           os.Getenv("DB_PORT"),
		DBUser:           os.Getenv("DB_USER"),
		DBName:           os.Getenv("DB_NAME"),
		DBPassword:       os.Getenv("DB_PASSWORD"),
		RedisAddr:        os.Getenv("REDIS_ADDR"),
		RedisPassword:    os.Getenv("REDIS_PASSWORD"),
		RedisDB:          redisDB,
		TelegramBotToken: os.Getenv("TELEGRAM_BOT_TOKEN"),
		Port:             os.Getenv("PORT"),
	}, nil
}
