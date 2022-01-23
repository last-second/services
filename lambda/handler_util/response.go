package handler_util

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambdacontext"
	trace "github.com/hans-m-song/go-stacktrace"
)

func headers() map[string]string {
	return map[string]string{
		"X-Lambda-Name": lambdacontext.FunctionName,
	}
}

func RespondWithError(err error) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 500,
		Headers:    headers(),
		Body:       trace.String(err),
	}
}

func RespondWithSuccess(body string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    headers(),
		Body:       body,
	}
}
