package config

import (
	trace "github.com/hans-m-song/go-stacktrace"
	"github.com/sirupsen/logrus"
)

var (
	ErrorParseLevel = trace.New("ErrorParseLevel")
)

func initLogger() *trace.Error {
	level, err := logrus.ParseLevel(Values.Loglevel)
	if err != nil {
		return ErrorParseLevel.Trace(err).Add("level", Values.Loglevel)
	}

	logrus.SetLevel(level)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	return nil
}
