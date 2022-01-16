package config

import "github.com/sirupsen/logrus"

func Init() {
	if err := initConfig(); err != nil {
		logrus.WithError(err).Fatal("error initializing config")
	}

	if err := initLogger(); err != nil {
		logrus.WithError(err).Fatal("error initializing logger")
	}
}
