package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	Port      string         `mapstructure:"port"`
	JWTSecret string         `mapstructure:"jwt_secret"`
	Database  DBConfig       `mapstructure:"database"`
	Timeouts  ServerTimeouts `mapstructure:"timeouts"`
	Kafka     KafkaConfig    `mapstructure:"kafka"`
}

type DBConfig struct {
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
	Host     string `mapstructure:"host"`
	Port     string `mapstructure:"port"`
	DBName   string `mapstructure:"dbname"`
}

type ServerTimeouts struct {
	ReadTimeout  int `mapstructure:"read"`
	WriteTimeout int `mapstructure:"write"`
	IdleTimeout  int `mapstructure:"idle"`
}

type KafkaConfig struct {
	Providers string `mapstrusture:"providers"`
}

func MustLoadConfig(filepath string) *Config {
	viper.SetConfigFile(filepath)
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("error reading config file: %v", err))
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		panic(fmt.Errorf("error unmarshalling config: %v", err))
	}

	return &config
}
