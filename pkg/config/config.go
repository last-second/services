package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	// deployment

	Loglevel string

	// runtime

	Region string
	Stage  string
}

var Values *Config = nil

// Retrieves or instantiates singleton instance of config values
func initConfig() error {
	viper.AutomaticEnv()
	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		return err
	}

	values := &Config{}
	if err := viper.Unmarshal(values); err != nil {
		return err
	}
	Values = values

	return nil
}
