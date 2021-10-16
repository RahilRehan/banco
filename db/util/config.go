package util

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	DB_USER               string        `mapstructure:"DB_USER"`
	DB_NAME               string        `mapstructure:"DB_NAME"`
	DB_PORT               string        `mapstructure:"DB_PORT"`
	DB_HOST               string        `mapstructure:"DB_HOST"`
	DB_PASSWORD           string        `mapstructure:"DB_PASSWORD"`
	MIGRATIONS_PATH       string        `mapstructure:"MIGRATIONS_PATH"`
	DRIVER_NAME           string        `mapstructure:"DRIVER_NAME"`
	SSL_MODE              string        `mapstructure:"SSL_MODE"`
	TIMEOUT               string        `mapstructure:"TIMEOUT"`
	SERVER_ADDRESS        string        `mapstructure:"SERVER_ADDRESS"`
	ACCESS_TOKEN_DURATION time.Duration `mapstructure:"ACCESS_TOKEN_DURATION"`
}

func LoadConfig(path string) (cfg *Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigName("app")
	viper.SetConfigType("env")

	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&cfg)

	return
}
