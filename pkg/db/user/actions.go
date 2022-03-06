package user

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	trace "github.com/hans-m-song/go-stacktrace"
	"github.com/last-second/services/pkg/db"
)

var (
	ErrorCreateUser = trace.New("ERROR_CREATE_USER")
	ErrorGetUser    = trace.New("ERROR_GET_USER")
	ErrorUpdateUser = trace.New("ERROR_UPDATE_USER")
)

func GetUser(ctx context.Context, tableName string, user *User) (*User, error) {
	client, err := db.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	if err != nil {
		return nil, err
	}

	item, err := user.MarshalAttributes(true)
	if err != nil {
		return nil, err
	}

	getItemInput := dynamodb.GetItemInput{
		TableName: &tableName,
		Key:       item,
	}

	response, err := client.GetItem(ctx, &getItemInput)
	if err != nil {
		return nil, ErrorGetUser.Trace(err).
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

func CreateUser(ctx context.Context, tableName string, user *User) (*User, error) {
	client, err := db.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	item, err := user.MarshalAttributes(false)
	if err != nil {
		return nil, err
	}

	putItemInput := dynamodb.PutItemInput{
		TableName: &tableName,
		Item:      item,
	}

	response, err := client.PutItem(ctx, &putItemInput)
	if err != nil {
		return nil, ErrorCreateUser.Trace(err).
			Add("tableName", *putItemInput.TableName).
			Add("item", putItemInput.Item)
	}

	result, err := UnmarshalAttributes(response.Attributes)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func UpdateUser(ctx context.Context, tableName string, user *User) (*User, error) {
	client, err := db.GetClient(ctx)
	if err != nil {
		return nil, err
	}

	updateExpression := "SET updated_at=:updated_at"
	updateValues := map[string]types.AttributeValue{":updated_at": &types.AttributeValueMemberS{Value: time.Now().UTC().Format(time.RFC3339)}}

	if user.Email != "" {
		updateExpression += ", email=:email"
		updateValues[":email"] = &types.AttributeValueMemberS{Value: user.Email}
	}

	if user.UserName != "" {
		updateExpression += ", user_name=:user_name"
		updateValues[":user_name"] = &types.AttributeValueMemberS{Value: user.UserName}
	}

	updateItemInput := dynamodb.UpdateItemInput{
		TableName:                 &tableName,
		Key:                       map[string]types.AttributeValue{"id": &types.AttributeValueMemberS{Value: user.Id}},
		ReturnValues:              types.ReturnValueAllNew,
		UpdateExpression:          &updateExpression,
		ExpressionAttributeValues: updateValues,
	}

	response, err := client.UpdateItem(ctx, &updateItemInput)
	if err != nil {
		return nil, ErrorUpdateUser.Trace(err).
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
