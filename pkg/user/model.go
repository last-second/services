package user

import (
	"encoding/json"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	trace "github.com/hans-m-song/go-stacktrace"
	"github.com/last-second/services/pkg/db"

	"github.com/google/uuid"
)

var (
	ErrorUnmarshalUserAttributes = trace.New("ERROR_UNMARSHAL_USER_ATTRIBUTES")
	ErrorMarshalUserAttributes   = trace.New("ERROR_MARSHAL_USER_ATTRIBUTES")
	ErrorUnmarshalUser           = trace.New("ERROR_UNMARSHAL_USER")
	ErrorMarshalUser             = trace.New("ERROR_MARSHAL_USER")
	ErrorMissingFields           = trace.New("ERROR_MISSING_FIELDS")
)

type User struct {
	Id        string `json:"id"         dynamodbav:"id,         omitempty"`
	Email     string `json:"email"      dynamodbav:"email,      omitempty"`
	UserName  string `json:"user_name"  dynamodbav:"user_name,  omitempty"`
	CreatedAt string `json:"created_at" dynamodbav:"created_at, omitempty"`
	UpdatedAt string `json:"updated_at" dynamodbav:"updated_at, omitempty"`
}

// creates a new user with values set for id, createdAt, and updatedAt
func NewEmptyUser() *User {
	now := time.Now().UTC().Format(time.RFC3339)
	id := uuid.NewString()

	return &User{
		Id:        id,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func (user *User) MarshalAttributes(omitempty bool) (map[string]types.AttributeValue, error) {
	attrs, err := attributevalue.MarshalMap(user)
	if err != nil {
		return nil, ErrorMarshalUser.Trace(err).Add("user", user)
	}

	if !omitempty {
		return attrs, nil
	}

	return db.FilterDynamodbAttributevalueMap(attrs), nil
}

// creates a new user with the given values,
//
// also includes values from `NewEmptyUser`
func NewUser(email, userName string) *User {
	user := NewEmptyUser()
	user.Email = email
	user.UserName = userName
	return user
}

func UnmarshalAttributes(attributes map[string]types.AttributeValue) (*User, error) {
	result := User{}
	err := attributevalue.UnmarshalMap(attributes, &result)
	if err != nil {
		return nil, ErrorUnmarshalUserAttributes.Trace(err).Add("attributes", attributes)
	}

	return &result, nil
}

// unmarshals a raw value into a user and checks for required values
func Parse(raw []byte) (*User, error) {
	parsed := User{}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, ErrorUnmarshalUserAttributes.Trace(err).Add("raw", raw)
	}

	missing := []string{}

	if parsed.Email == "" {
		missing = append(missing, "Email")
	}

	if parsed.UserName == "" {
		missing = append(missing, "UserName")
	}

	if len(missing) > 0 {
		return nil, ErrorMissingFields.Tracef("missing required field(s)").Add("fields", missing)
	}

	return &parsed, nil
}

func FromMap(values map[string]string) (*User, error) {
	serialised, err := json.Marshal(values)
	if err != nil {
		return nil, ErrorMarshalUser.Trace(err).Add("values", values)
	}

	result := User{}
	json.Unmarshal(serialised, &result)
	return &result, nil
}
