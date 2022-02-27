package user

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	trace "github.com/hans-m-song/go-stacktrace"
	"github.com/last-second/services/pkg/db"
)

var (
	ErrorPutUser    = trace.New("ErrorPutUser")
	ErrorGetUser    = trace.New("ErrorGetUser")
	ErrorListUser   = trace.New("ErrorListUser")
	ErrorUpdateUser = trace.New("ErrorUpdateUser")
	ErrorDeleteUser = trace.New("ErrorDeleteUser")
)

func ListUsers() {
	// TODO
}

func GetUser(tableName string, user *User) (*User, error) {
	client, err := db.GetClient()
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

	response, err := client.GetItem(context.TODO(), &getItemInput)
	if err != nil {
		return nil, ErrorGetUser.Trace(err).Add("getItemInput", getItemInput)
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

func CreateUser(tableName string, user *User) (*User, error) {
	client, err := db.GetClient()
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

	response, err := client.PutItem(context.TODO(), &putItemInput)
	if err != nil {
		return nil, ErrorPutUser.Trace(err).Add("putItemInput", putItemInput)
	}

	result, err := UnmarshalAttributes(response.Attributes)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func UpdateUser() {
	// TODO
}

func DeleteUser() {
	// TODO
}
