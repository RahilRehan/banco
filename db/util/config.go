package util

import (
	"github.com/spf13/viper"
)

type Config struct {
	DB_USER        string `mapstructure:"DB_USER"`
	DB_NAME        string `mapstructure:"DB_NAME"`
	DB_PORT        string `mapstructure:"DB_PORT"`
	DB_HOST        string `mapstructure:"DB_HOST"`
	DRIVER_NAME    string `mapstructure:"DRIVER_NAME"`
	SSL_MODE       string `mapstructure:"SSL_MODE"`
	TIMEOUT        string `mapstructure:"TIMEOUT"`
	SERVER_ADDRESS string `mapstructure:"SERVER_ADDRESS"`
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
