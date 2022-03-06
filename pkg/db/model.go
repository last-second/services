package db

import (
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	trace "github.com/hans-m-song/go-stacktrace"
)

type Action int

const (
	GetAction Action = iota
	CreateAction
	UpdateAction
	DeleteAction
)

var (
	ErrorMissingAttributes   = trace.New("ERROR_MISSING_MODEL_ATTRIBUTES")
	ErrorForbiddenAttributes = trace.New("ERROR_FORBIDDEN_MODEL_ATTRIBUTES")
	ErrorInvalidAttributes   = trace.New("ERROR_INVALID_MODEL_ATTRIBUTES")

	ErrorUnmarshalModelAttributes = trace.New("ERROR_UNMARSHAL_MODEL_ATTRIBUTES")
	ErrorUnmarshalModel           = trace.New("ERROR_UNMARSHAL_MODEL_ATTRIBUTES")

	ErrorMarshalModelAttributes = trace.New("ERROR_UNMARSHAL_MODEL")
	ErrorMarshalModel           = trace.New("ERROR_UNMARSHAL_MODEL")
)

type Model interface {
	Marshal() ([]byte, error)
	MarshalAttributes(omitempty bool) (map[string]types.AttributeValue, error)
	EnsureAttributes(action Action) error
}

func ContainsAction(actions []Action, action Action) bool {
	for _, v := range actions {
		if action == v {
			return true
		}
	}

	return false
}
