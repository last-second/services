package handler_util

import (
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	trace "github.com/hans-m-song/go-stacktrace"
	"github.com/last-second/services/pkg/api"
	"github.com/sirupsen/logrus"
)

func headers() map[string]string {
	return map[string]string{
		"X-Lambda-Name": lambdacontext.FunctionName,
		"Content-Type":  "application/json",
	}
}

func RespondWithError(code int, err error, message string) (events.APIGatewayProxyResponse, error) {
	traced := trace.Guarantee(err).Tracef(message)
	logrus.WithField("error", traced).Error("request failed")
	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Headers:    headers(),
		Body:       api.MustMarshal(api.RespondWithError(traced.Name, message)),
	}, err
}

func RespondWithSuccess(data interface{}) (events.APIGatewayProxyResponse, error) {
	logrus.WithField("data", data).Error("request succeeded")
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Headers:    headers(),
		Body:       api.MustMarshal(api.RespondWithSuccess(data)),
	}, nil
}
