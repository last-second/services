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
	"github.com/last-second/services/pkg/db/task"
	"github.com/sirupsen/logrus"
)

func init() {
	logrus.Info("Initializing GetTask")
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
	logrus.WithFields(handler_util.EventMeta(event)).Debug("Executing GetTask")

	id, ok := event.QueryStringParameters["id"]
	if !ok {
		return handler_util.RespondWithError(http.StatusBadRequest, api.ErrorInvalidQuery, "A task id is required")
	}

	foundTask, err := task.GetTask(ctx, config.Values.TasktableName, &task.Task{Id: id})
	if err != nil {
		return handler_util.RespondWithError(http.StatusBadRequest, err, "Could not get task")
	}

	if foundTask == nil {
		return handler_util.RespondWithError(http.StatusNotFound, api.ErrorNotFound, "Could not find matching task")
	}

	return handler_util.RespondWithSuccess(foundTask)
}

func main() {
	lambda.Start(handler)
}
