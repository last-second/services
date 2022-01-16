package config

import (
	"github.com/sirupsen/logrus"
)

func initLogger() error {
	level, err := logrus.ParseLevel(Values.Loglevel)
	if err != nil {
		return err
	}

	logrus.SetLevel(level)
	return nil
}
