package test

import (
	"fmt"
	"testing"
)

type Assert struct {
	T *testing.T
}

func (a Assert) ErrorIsNil(err error) {
	if err != nil {
		a.T.Fatalf("received: \"%s\", wanted nil error", err.Error())
	}
}

func (a Assert) MatchesString(received, wanted string) {
	if wanted != received {
		a.T.Fatalf("received: \"%s\", wanted: \"%s\"", received, wanted)
	}
}

func (a Assert) MatchesInt(received, wanted int) {
	if wanted != received {
		a.T.Fatalf("received: \"%d\", wanted: \"%d\"", received, wanted)
	}
}

func (a Assert) MapHasKey(value map[string]interface{}, key string) {
	if _, ok := value[key]; !ok {
		a.T.Fatalf("wanted key did not exist: \"%s\"", key)
	}
}

func (a Assert) IsString(value interface{}) {
	if fmt.Sprintf("%T", value) != "string" {
		a.T.Fatalf("received: \"%T\", wanted type \"string\"", value)
	}
}
