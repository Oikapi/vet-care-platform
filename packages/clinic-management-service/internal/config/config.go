package config

import "github.com/spf13/viper"

type Config struct {
    MySQLDSN    string
    RedisAddr   string
}

func LoadConfig() (*Config, error) {
    viper.SetConfigName("config")
    viper.AddConfigPath(".")
    viper.SetConfigType("yaml")

    if err := viper.ReadInConfig(); err != nil {
        return nil, err
    }

    return &Config{
        MySQLDSN:  viper.GetString("mysql.dsn"),
        RedisAddr: viper.GetString("redis.addr"),
    }, nil
}