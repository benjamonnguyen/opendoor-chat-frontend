package config

import (
	"log"

	"github.com/spf13/viper"
)

type Config struct {
	Salt           string
	BackendBaseUrl string `mapstructure:"backend_base_url"`
}

func LoadConfig(file string) Config {
	viper.AutomaticEnv()
	viper.SetConfigFile(file)
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal(err)
	}

	return cfg
}
