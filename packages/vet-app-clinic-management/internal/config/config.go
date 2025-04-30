package config

import (
	"os"
	"fmt"
)

type Config struct {
	MySQLDSN    string
	RedisAddr   string
    SwaggerHost string `mapstructure:"SWAGGER_HOST"`
}

func LoadConfig() (*Config, error) {
	cfg := &Config{
		MySQLDSN:  "root:password@tcp(localhost:3306)/clinic_db?charset=utf8mb4&parseTime=True&loc=Local",
		RedisAddr: "localhost:6379",
        SwaggerHost: "localhost:8080",
	}

	if dsn := os.Getenv("MYSQL_DSN"); dsn != "" {
		cfg.MySQLDSN = dsn
	}

	if addr := os.Getenv("REDIS_ADDR"); addr != "" {
		cfg.RedisAddr = addr
	}

	if cfg.MySQLDSN == "" {
		return nil, fmt.Errorf("mysql DSN configuration is required")
	}



	return cfg, nil
}