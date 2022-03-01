package task

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	trace "github.com/hans-m-song/go-stacktrace"
	"github.com/last-second/services/pkg/db"
)

type FrequencyType string

var (
	ErrorUnmarshalTaskAttributes = trace.New("ERROR_UNMARSHAL_TASK_ATTRIBUTES")
	ErrorMarshalTaskAttributes   = trace.New("ERROR_MARSHAL_TASK_ATTRIBUTES")
	ErrorUnmarshalTask           = trace.New("ERROR_UNMARSHAL_TASK")
	ErrorMarshalTask             = trace.New("ERROR_MARSHAL_TASK")
	ErrorTaskFields              = trace.New("ERROR_TASK_FIELDS")
	ErrorTaskFrequency           = trace.New("ERROR_TASK_FREQUENCY")
)

const (
	FrequencyTypeNever FrequencyType = "NEVER"
	FrequencyTypeDay   FrequencyType = "DAY"
	FrequencyTypeWeek  FrequencyType = "WEEK"
	FrequencyTypeMonth FrequencyType = "MONTH"
	FrequencyTypeYear  FrequencyType = "YEAR"
)

type Task struct {
	Id            string        `json:"id"             dynamodbav:"id,             omitempty"`
	UserId        string        `json:"user_id"        dynamodbav:"user_id,        omitempty"`
	Title         string        `json:"title"          dynamodbav:"title,          omitempty"`
	Description   string        `json:"description"    dynamodbav:"description,    omitempty"`
	StartAt       string        `json:"start_at"       dynamodbav:"start_at,       omitempty"`
	EndAt         string        `json:"end_at"         dynamodbav:"end_at,         omitempty"`
	FrequencyType FrequencyType `json:"frequency_type" dynamodbav:"frequency_type, omitempty"`
	Frequency     int           `json:"frequency"      dynamodbav:"frequency,      omitempty"`
	CreatedAt     string        `json:"created_at"     dynamodbav:"created_at,     omitempty"`
	UpdatedAt     string        `json:"updated_at"     dynamodbav:"updated_at,     omitempty"`
}

func NewEmptyTask() *Task {
	now := time.Now().UTC().Format(time.RFC3339)
	id := uuid.NewString()

	return &Task{
		Id:            id,
		CreatedAt:     now,
		UpdatedAt:     now,
		FrequencyType: FrequencyTypeNever,
	}
}

func NewTask(
	userId, title, description, startAt, endAt string,
	frequencyType FrequencyType,
	frequency int,
) *Task {
	task := NewEmptyTask()
	task.UserId = userId
	task.Title = title
	task.Description = description
	task.StartAt = startAt
	task.EndAt = endAt
	task.FrequencyType = frequencyType
	task.Frequency = frequency
	return task
}

func (task *Task) MarshalAttributes(omitempty bool) (map[string]types.AttributeValue, error) {
	attrs, err := attributevalue.MarshalMap(task)
	if err != nil {
		return nil, ErrorMarshalTaskAttributes.Trace(err).Add("task", task)
	}

	if !omitempty {
		return attrs, nil
	}

	return db.FilterDynamodbAttributevalueMap(attrs), nil
}

func UnmarshalAttributes(attributes map[string]types.AttributeValue) (*Task, error) {
	result := Task{}
	if err := attributevalue.UnmarshalMap(attributes, &result); err != nil {
		return nil, ErrorUnmarshalTaskAttributes.Trace(err).Add("attributes", attributes)
	}

	return &result, nil
}

func (task *Task) EnsureCreationAttributes() error {
	missing := []string{}

	if task.UserId == "" {
		missing = append(missing, "UserId")
	}

	if task.Title == "" {
		missing = append(missing, "Title")
	}

	if task.StartAt == "" {
		missing = append(missing, "StartAt")
	}

	if task.FrequencyType == "" {
		missing = append(missing, "FrequencyType")
	}

	if task.FrequencyType != FrequencyTypeNever && task.Frequency == 0 {
		missing = append(missing, "Frequency")
	}

	if task.FrequencyType != FrequencyTypeNever && task.EndAt == "" {
		missing = append(missing, "EndAt")
	}

	if len(missing) > 0 {
		return ErrorTaskFields.Tracef("missing required field(s)").Add("fields", missing)
	}

	disallowed := []string{}

	if task.Id != "" {
		disallowed = append(disallowed, "Id")
	}

	if task.UserId != "" {
		disallowed = append(disallowed, "UserId")
	}

	if task.CreatedAt != "" {
		disallowed = append(disallowed, "CreatedAt")
	}

	if task.UpdatedAt != "" {
		disallowed = append(disallowed, "UpdatedAt")
	}

	if len(disallowed) > 0 {
		return ErrorTaskFields.Tracef("must not specify field(s)").Add("fields", disallowed)
	}

	return nil
}

func (task *Task) ValidateFrequency() error {
	if task.FrequencyType == FrequencyTypeNever {
		return nil
	}

	if task.Frequency == 0 {
		return ErrorTaskFrequency.Tracef("\"frequency\" must be set if \"frequency_type\" is \"%s\"", FrequencyTypeNever)
	}

	if task.EndAt == "" {
		return ErrorTaskFrequency.Tracef("\"end_at\" must be set if \"frequency_type\" is \"%s\"", FrequencyTypeNever)
	}

	return nil
}
