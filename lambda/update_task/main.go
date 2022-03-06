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
	logrus.WithFields(handler_util.EventMeta(event)).Debug("Executing UpdateTask")

	partialTask := task.Task{}
	if err := json.Unmarshal([]byte(event.Body), &partialTask); err != nil {
		return handler_util.RespondWithError(http.StatusBadRequest, api.ErrorInvalidBody.Add("body", event.Body).Trace(err), "Could not parse body")
	}

	if partialTask.Id == "" || partialTask.CreatedAt != "" || partialTask.UpdatedAt != "" {
		return handler_util.RespondWithError(http.StatusBadRequest, api.ErrorInvalidBody.Add("task", partialTask), "Must specify Id and optionally email and task_name when updating a task")
	}

	updatedTask, err := task.UpdateTask(ctx, config.Values.TasktableName, &partialTask)
	if err != nil {
		return handler_util.RespondWithError(http.StatusInternalServerError, trace.Guarantee(err).Add("task", partialTask), "Error updating task")
	}

	return handler_util.RespondWithSuccess(updatedTask)
}

func main() {
	lambda.Start(handler)
}
