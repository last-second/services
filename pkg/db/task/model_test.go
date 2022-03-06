package task_test

import (
	"testing"

	trace "github.com/hans-m-song/go-stacktrace"
	"github.com/last-second/services/pkg/db"
	"github.com/last-second/services/pkg/db/task"
	"github.com/last-second/services/pkg/test"
)

func TestEnsureFieldsEmptyTask(t *testing.T) {
	mockTask := task.Task{
		Id:            "",
		UserId:        "",
		Title:         "",
		Description:   "",
		StartAt:       "",
		EndAt:         "",
		Frequency:     0,
		FrequencyType: "",
		CreatedAt:     "",
		UpdatedAt:     "",
	}

	err := mockTask.EnsureAttributes(db.CreateAction)
	meta := trace.Guarantee(err).Meta

	assert := test.Assert{T: t}
	assert.MapHasKey(meta, "fields")
	assert.IsString(meta["fields"])
	assert.MatchesString(meta["fields"].(string), "[UserId Title StartAt FrequencyType Frequency EndAt]")
}

func TestEnsureFieldsPartialTask(t *testing.T) {
	mockTask := task.Task{
		Id:            "",
		UserId:        "",
		Title:         "title",
		Description:   "",
		StartAt:       "",
		EndAt:         "",
		Frequency:     0,
		FrequencyType: "",
		CreatedAt:     "",
		UpdatedAt:     "",
	}

	err := mockTask.EnsureAttributes(db.CreateAction)
	meta := trace.Guarantee(err).Meta

	assert := test.Assert{T: t}
	assert.MapHasKey(meta, "fields")
	assert.IsString(meta["fields"])
	assert.MatchesString(meta["fields"].(string), "[UserId StartAt FrequencyType Frequency EndAt]")
}

func TestEnsureFieldsFrequencyTypeIsNotNever(t *testing.T) {
	mockTask := task.Task{
		Id:            "",
		UserId:        "user id",
		Title:         "title",
		Description:   "description",
		StartAt:       "start at",
		EndAt:         "",
		Frequency:     0,
		FrequencyType: task.FrequencyTypeDay,
		CreatedAt:     "",
		UpdatedAt:     "",
	}

	err := mockTask.EnsureAttributes(db.CreateAction)
	meta := trace.Guarantee(err).Meta

	assert := test.Assert{T: t}
	assert.MapHasKey(meta, "fields")
	assert.IsString(meta["fields"])
	assert.MatchesString(meta["fields"].(string), "[Frequency EndAt]")

}
func TestEnsureFieldsFrequencyTypeIsNever(t *testing.T) {
	mockTask := task.Task{
		Id:            "",
		UserId:        "user id",
		Title:         "title",
		Description:   "description",
		StartAt:       "start at",
		EndAt:         "",
		Frequency:     0,
		FrequencyType: task.FrequencyTypeNever,
		CreatedAt:     "",
		UpdatedAt:     "",
	}

	err := mockTask.EnsureAttributes(db.CreateAction)

	assert := test.Assert{T: t}
	assert.ErrorIsNil(err)
}
