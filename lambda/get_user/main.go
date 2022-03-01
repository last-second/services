package main

import (
	"context"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	trace "github.com/hans-m-song/go-stacktrace"
	"github.com/last-second/services/lambda/handler_util"
	"github.com/last-second/services/pkg/api"
	"github.com/last-second/services/pkg/config"
	"github.com/last-second/services/pkg/db"
	"github.com/last-second/services/pkg/user"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.Info("Initializing GetUser")
	config.Init()

	if _, err := db.GetClient(context.Background()); err != nil {
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
	logrus.WithFields(handler_util.EventMeta(event)).Debug("Executing GetUser")

	id, ok := event.QueryStringParameters["id"]
	if !ok {
		return handler_util.RespondWithError(http.StatusBadRequest, api.ErrorInvalidQuery, "A user id is required")
	}

	foundUser, err := user.GetUser(ctx, config.Values.UsertableName, &user.User{Id: id})
	if err != nil {
		return handler_util.RespondWithError(http.StatusBadRequest, err, "Could not get user")
	}

	if foundUser == nil {
		return handler_util.RespondWithError(http.StatusNotFound, err, "Could not find matching user")
	}

	return handler_util.RespondWithSuccess(foundUser)
}

func main() {
	lambda.Start(handler)
}
