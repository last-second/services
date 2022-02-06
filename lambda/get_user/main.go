package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	trace "github.com/hans-m-song/go-stacktrace"
	"github.com/last-second/services/lambda/handler_util"
	"github.com/last-second/services/pkg/config"
	"github.com/last-second/services/pkg/db"
	"github.com/last-second/services/pkg/user"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.Info("Initializing GetUser")
	config.Init()

	if _, err := db.GetClient(); err != nil {
		logrus.Fatal(trace.String(err))
	}
}

func handler(
	ctx context.Context,
	event events.APIGatewayProxyRequest,
) (
	events.APIGatewayProxyResponse,
	error,
) {
	logrus.Info("Executing GetUser")
	logrus.Debugf("event: %+v", event)

	search, err := user.FromMap(event.QueryStringParameters)
	if err != nil {
		return handler_util.RespondWithError(err), err
	}

	logrus.Debugf("query: %+v", search)
	resp, err := user.GetUser(config.Values.UsertableName, search)
	if err != nil {
		return handler_util.RespondWithError(err), err
	}

	serialised, err := json.Marshal(resp)
	if err != nil {
		return handler_util.RespondWithError(err), err
	}

	logrus.Debugf("%+v", resp)
	return handler_util.RespondWithSuccess(string(serialised)), nil
}

func main() {
	lambda.Start(handler)
}
