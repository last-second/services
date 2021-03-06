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
	_ db.Model = (*User)(nil)

	ErrorUnmarshalUserAttributes = trace.New("ERROR_UNMARSHAL_USER_ATTRIBUTES")
	ErrorMarshalUserAttributes   = trace.New("ERROR_MARSHAL_USER_ATTRIBUTES")
	ErrorUnmarshalUser           = trace.New("ERROR_UNMARSHAL_USER")
	ErrorMarshalUser             = trace.New("ERROR_MARSHAL_USER")
	ErrorUserFields              = trace.New("ERROR_USER_FIELDS")
)

type User struct {
	Id        string `json:"id"         dynamodbav:"id,         omitempty"`
	Email     string `json:"email"      dynamodbav:"email,      omitempty"`
	UserName  string `json:"user_name"  dynamodbav:"user_name,  omitempty"`
	CreatedAt string `json:"created_at" dynamodbav:"created_at, omitempty"`
	UpdatedAt string `json:"updated_at" dynamodbav:"updated_at, omitempty"`
}

func (user *User) Marshal() ([]byte, error) {
	result, err := json.Marshal(user)
	if err != nil {
		return nil, db.ErrorMarshalModel.Trace(err)
	}

	return result, err
}

func (user *User) MarshalAttributes(omitempty bool) (map[string]types.AttributeValue, error) {
	attrs, err := attributevalue.MarshalMap(user)
	if err != nil {
		return nil, db.ErrorMarshalModelAttributes.Trace(err).Add("user", user)
	}

	if omitempty {
		return db.FilterDynamodbAttributevalueMap(attrs), nil
	}

	return attrs, nil
}

func (user *User) EnsureAttributes(action db.Action) error {
	required := []string{}

	if user.Email == "" {
		required = append(required, "Email")
	}

	if user.UserName == "" {
		required = append(required, "UserName")
	}

	if len(required) > 0 {
		return ErrorUserFields.Tracef("missing required field(s)").Add("fields", required)
	}

	disallowed := []string{}

	if user.Id != "" {
		disallowed = append(disallowed, "Id")
	}

	if user.CreatedAt != "" {
		disallowed = append(disallowed, "CreatedAt")
	}

	if user.UpdatedAt != "" {
		disallowed = append(disallowed, "UpdatedAt")
	}

	if len(disallowed) > 0 {
		return ErrorUserFields.Tracef("must not specify field(s)").Add("fields", disallowed)
	}

	return nil
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
func Unmarshal(raw []byte) (*User, error) {
	parsed := User{}
	if err := json.Unmarshal(raw, &parsed); err != nil {
		return nil, ErrorUnmarshalUser.Trace(err).Add("raw", raw)
	}

	return &parsed, nil
}
