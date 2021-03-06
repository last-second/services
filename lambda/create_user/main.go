package main

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	trace "github.com/hans-m-song/go-stacktrace"
	"github.com/last-second/services/lambda/handler_util"
	"github.com/last-second/services/pkg/api"
	"github.com/last-second/services/pkg/config"
	"github.com/last-second/services/pkg/db"
	"github.com/last-second/services/pkg/db/user"
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
	logrus.WithFields(handler_util.EventMeta(event)).Debug("Executing CreateUser")

	partialUser := user.User{}
	if err := json.Unmarshal([]byte(event.Body), &partialUser); err != nil {
		return handler_util.RespondWithError(http.StatusBadRequest, api.ErrorInvalidBody.Add("body", event.Body).Trace(err), "Could not parse body")
	}

	if err := partialUser.EnsureAttributes(db.CreateAction); err != nil {
		return handler_util.RespondWithError(http.StatusBadRequest, err, "Can only specify email and user_name when creating a user")
	}

	pendingUser := user.NewUser(partialUser.Email, partialUser.UserName)
	createdUser, err := user.CreateUser(ctx, config.Values.UsertableName, pendingUser)
	if err != nil {
		return handler_util.RespondWithError(http.StatusInternalServerError, trace.Guarantee(err).Add("user", partialUser), "Error creating user")
	}

	return handler_util.RespondWithSuccess(createdUser)
}

func main() {
	lambda.Start(handler)
}
