package task

import (
	"context"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	trace "github.com/hans-m-song/go-stacktrace"
	"github.com/last-second/services/pkg/db"
)

var (
	ErrorCreateTask = trace.New("ERROR_CREATE_USER")
	ErrorGetTask    = trace.New("ERROR_GET_USER")
	ErrorUpdateTask = trace.New("ERROR_UPDATE_USER")
)

func GetTask(ctx context.Context, tableName string, task *Task) (*Task, error) {
	client, err := db.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	item, err := task.MarshalAttributes(true)
	if err != nil {
		return nil, err
	}

	getItemInput := dynamodb.GetItemInput{
		TableName: &tableName,
		Key:       item,
	}

	response, err := client.GetItem(ctx, &getItemInput)
	if err != nil {
		return nil, ErrorGetTask.Trace(err).
			Add("tableName", *getItemInput.TableName).
			Add("key", getItemInput.Key)
	}

	if len(response.Item) < 1 {
		return nil, nil
	}

	result, err := UnmarshalAttributes(response.Item)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func CreateTask(ctx context.Context, tableName string, task *Task) (*Task, error) {
	client, err := db.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	item, err := task.MarshalAttributes(false)
	if err != nil {
		return nil, err
	}

	putItemInput := dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      item,
	}

	response, err := client.PutItem(ctx, &putItemInput)
	if err != nil {
		return nil, ErrorCreateTask.Trace(err).
			Add("tableName", *putItemInput.TableName).
			Add("item", putItemInput.Item)
	}

	result, err := UnmarshalAttributes(response.Attributes)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func UpdateTask(ctx context.Context, tableName string, task *Task) (*Task, error) {
	client, err := db.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	if err := task.ValidateFrequency(); err != nil {
		return nil, err
	}

	updateExpression := "SET updated_at=:updated_at"
	updateValues := map[string]types.AttributeValue{":updated_at": &types.AttributeValueMemberS{Value: time.Now().UTC().Format(time.RFC3339)}}

	if task.Title != "" {
		updateExpression += "title=:title"
		updateValues[":title"] = &types.AttributeValueMemberS{Value: task.Title}
	}

	if task.Description != "" {
		updateExpression += "description=:description"
		updateValues[":description"] = &types.AttributeValueMemberS{Value: task.Description}
	}

	if task.StartAt != "" {
		updateExpression += "start_at=:start_at"
		updateValues[":start_at"] = &types.AttributeValueMemberS{Value: task.StartAt}
	}

	if task.EndAt != "" {
		updateExpression += "end_at=:end_at"
		updateValues[":end_at"] = &types.AttributeValueMemberS{Value: task.EndAt}
	}

	if task.FrequencyType != "" {
		updateExpression += "frequency_type=:frequency_type"
		updateValues[":frequency_type"] = &types.AttributeValueMemberS{Value: string(task.FrequencyType)}
	}

	if task.FrequencyType != FrequencyTypeNever && task.Frequency > 0 {
		updateExpression += "frequency=:frequency"
		updateValues[":frequency"] = &types.AttributeValueMemberN{Value: strconv.Itoa(task.Frequency)}
	}

	updateItemInput := dynamodb.UpdateItemInput{
		TableName:                 &tableName,
		Key:                       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: task.Id}},
		ReturnValues:              types.ReturnValueAllNew,
		UpdateExpression:          &updateExpression,
		ExpressionAttributeValues: updateValues,
	}

	response, err := client.UpdateItem(ctx, &updateItemInput)
	if err != nil {
		return nil, ErrorUpdateTask.Trace(err).
			Add("tableName", *updateItemInput.TableName).
			Add("key", updateItemInput.Key).
			Add("updateExpression", *updateItemInput.UpdateExpression).
			Add("updateValues", updateValues)
	}

	result, err := UnmarshalAttributes(response.Attributes)
	if err != nil {
		return nil, err
	}

	return result, nil
}
