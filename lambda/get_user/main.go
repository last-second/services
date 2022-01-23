package main

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	trace "github.com/hans-m-song/go-stacktrace"
	"github.com/last-second/services/lambda/handler_util"
	"github.com/last-second/services/pkg/config"
	"github.com/last-second/services/pkg/db"
	"github.com/last-second/services/pkg/user"
	"github.com/sirupsen/logrus"
)

var (
	dbClient *dynamodb.Client
)

func init() {
	logrus.Info("Initializing GetUser")
	config.Init()

	client, err := db.GetClient()
	if err != nil {
		logrus.Fatal(trace.String(err))
	}

	dbClient = client
}

func handler(ctx context.Context,
	event events.APIGatewayProxyRequest,
) (
	response events.APIGatewayProxyResponse,
	err error,
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
