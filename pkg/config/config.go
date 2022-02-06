package config

import (
	"os"

	trace "github.com/hans-m-song/go-stacktrace"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	Values *Config

	ErrorReadConfig      = trace.New("ErrorReadConfig")
	ErrorUnmarshalConfig = trace.New("ErrorUnmarshalConfig")

	envDefaults = map[string]string{
		"LOGLEVEL":       "debug",
		"STAGE":          "dev",
		"AWS_PROFILE":    "",
		"AWS_REGION":     "",
		"USERTABLE_NAME": "",
	}
)

func setFromEnvOrDefault(key, defaultValue string) {
	if value, ok := os.LookupEnv(key); ok && value != "" {
		viper.SetDefault(key, value)
	} else {
		viper.SetDefault(key, defaultValue)
	}
}

type Config struct {
	// deployment

	Loglevel string `mapstructure:"loglevel"`

	// runtime

	Stage         string `mapstructure:"stage"`
	AwsProfile    string `mapstructure:"aws_profile"`
	AwsRegion     string `mapstructure:"aws_region"`
	UsertableName string `mapstructure:"usertable_name"`
}

// Retrieves or instantiates singleton instance of config values
func initConfig() *trace.Error {
	viper.AddConfigPath("../..") // for handlers
	viper.SetConfigName("config")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			logrus.Debug("no config file found")
			return ErrorReadConfig.Trace(err)
		}
	}

	for key, defaultValue := range envDefaults {
		setFromEnvOrDefault(key, defaultValue)
	}

	if Values == nil {
		Values = &Config{}
		if err := viper.Unmarshal(Values); err != nil {
			return ErrorUnmarshalConfig.Trace(err).Add("config", viper.AllSettings())
		}
	}

	return nil
}
