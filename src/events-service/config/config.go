package config

import (
	"log"

	"github.com/spf13/viper"
)

func LoadConfig() *viper.Viper {
	vp := viper.New()
	vp.SetConfigName("config")
	vp.SetConfigType("json")
	vp.AddConfigPath("./src/events-service/config")

	// in case we are injecting config from src
	vp.AddConfigPath("../config")

	if err := vp.ReadInConfig(); err != nil {
		log.Fatalf("Error reading in config file: %v", err)
	}
	return vp
}
