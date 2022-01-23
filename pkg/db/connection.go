package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	trace "github.com/hans-m-song/go-stacktrace"
	"github.com/sirupsen/logrus"
)

var (
	ErrorGettingConfig = trace.New("ErrorGettingConfig")
	client             *dynamodb.Client
)

// Singleton client getter
func GetClient() (*dynamodb.Client, error) {
	if client != nil {
		return client, nil
	}

	logrus.Info("getting aws sdk configuration")
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, ErrorGettingConfig.Trace(err)
	}

	logrus.Info("creating dynamodb client")
	client = dynamodb.NewFromConfig(cfg)

	return client, nil
}
