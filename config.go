package main

import (
	"log"

	"github.com/spf13/viper"
)

var cfg struct {
	Salt string
}

func init() {
	viper.AutomaticEnv()
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	if err := viper.Unmarshal(&cfg); err != nil {
		log.Fatal(err)
	}
}
