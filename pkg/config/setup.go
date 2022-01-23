package config

import (
	"os"

	trace "github.com/hans-m-song/go-stacktrace"
	"github.com/sirupsen/logrus"
)

func Init() {
	if err := initConfig(); err != nil {
		logrus.Fatal(trace.Guarantee(err).String())
	}

	logrus.Debugf("environment: %+v\n", os.Environ())
	logrus.Debugf("configuration: %+v\n", Values)

	if err := initLogger(); err != nil {
		logrus.Fatal(trace.Guarantee(err).String())
	}

}
