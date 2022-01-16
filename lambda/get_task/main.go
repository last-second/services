package main

import (
	"context"

	"github.com/aws/aws-lambda-go/lambda"
	h "github.com/last-second/services/lambda/handler"
	"github.com/last-second/services/pkg/config"
	"github.com/sirupsen/logrus"
)

type getTaskEvent struct {
}

type getTaskResponse struct {
	Raw string `json:"raw"`
}

func handler(ctx context.Context, event getTaskEvent) (h.Response, error) {
	logrus.Info(ctx, event)

	response := h.Response{
		Message: "GetTaskResponse",
		Body:    &getTaskResponse{Raw: "raw"},
	}

	return response, nil
}

func main() {
	config.Init()
	logrus.Info("Executing GetTask")
	lambda.Start(handler)
}
