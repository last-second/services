package handler_util

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/last-second/services/pkg/config"
	"github.com/sirupsen/logrus"
)

func EventMeta(event events.APIGatewayProxyRequest) logrus.Fields {
	return logrus.Fields{
		"SourceIP":   event.RequestContext.Identity.SourceIP,
		"HTTPMethod": event.HTTPMethod,
		"Path":       event.Path,
		"Stage":      config.Values.Stage,
	}
}
