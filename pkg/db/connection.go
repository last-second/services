package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	trace "github.com/hans-m-song/go-stacktrace"
	c "github.com/last-second/services/pkg/config"
	"github.com/sirupsen/logrus"
)

var (
	ErrorGettingConfig = trace.New("ErrorGettingConfig")
	client             *dynamodb.Client
)

var customResolver = aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
	if c.Values.Stage == "local" && service == dynamodb.ServiceID {
		return aws.Endpoint{
			PartitionID:   "aws",
			URL:           "http://dynamodb:8000",
			SigningRegion: c.Values.AwsRegion,
		}, nil
	}

	return aws.Endpoint{}, &aws.EndpointNotFoundError{}
})

// Singleton client getter
func GetClient() (*dynamodb.Client, error) {
	if client != nil {
		return client, nil
	}

	logrus.Info("getting aws sdk configuration")
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
	)

	if err != nil {
		return nil, ErrorGettingConfig.Trace(err)
	}

	logrus.Info("creating dynamodb client")
	client = dynamodb.NewFromConfig(cfg)

	return client, nil
}
