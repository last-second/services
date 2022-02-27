package config

import (
	"fmt"

	trace "github.com/hans-m-song/go-stacktrace"
	"github.com/sirupsen/logrus"
)

var (
	ErrorParseLevel      = trace.New("ErrorParseLevel")
	defaultJsonFormatter = logrus.JSONFormatter{}
	defaultTextFormatter = logrus.TextFormatter{
		ForceQuote:             true,
		QuoteEmptyFields:       true,
		DisableLevelTruncation: true,
	}
)

type jsonFormatter struct{}

func (f jsonFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	for k, v := range entry.Data {
		if err, ok := v.(*trace.Error); ok {
			// extract error details for json formatting otherwise logrus calls Error() instead
			entry.Data[k] = map[string]interface{}{
				"name":    err.Name,
				"message": err.Message,
				"meta":    err.Meta,
				"stack":   err.Stack,
			}
		}
	}

	return defaultJsonFormatter.Format(entry)
}

type textFormatter struct{}

func (f textFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	data := make(logrus.Fields)
	errStr := ""

	// override data to skip formatting custom error
	for k, v := range entry.Data {
		switch v := v.(type) {
		case trace.SerializableError:
			if errStr != "" {
				errStr += "\n"
			}
			errStr += fmt.Sprintf("%s=\"%s\"", k, v.String())
		default:
			data[k] = v
		}
	}

	entry.Data = data
	rest, err := defaultTextFormatter.Format(entry)
	if errStr != "" {
		rest = append(rest, errStr...)
	}
	return rest, err
}

func initLogger() *trace.Error {
	level, err := logrus.ParseLevel(Values.Loglevel)
	if err != nil {
		return ErrorParseLevel.Trace(err).Add("level", Values.Loglevel)
	}

	logrus.SetLevel(level)

	if Values.Stage == "local" {
		logrus.SetFormatter(new(textFormatter))
	} else {
		logrus.SetFormatter(new(jsonFormatter))
	}

	return nil
}
