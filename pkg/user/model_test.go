package user_test

import (
	"testing"

	trace "github.com/hans-m-song/go-stacktrace"
	"github.com/last-second/services/pkg/test"
	"github.com/last-second/services/pkg/user"
)

func TestEnsureCreationAttributesEmptyUser(t *testing.T) {
	mockUser := user.User{
		Id:        "",
		Email:     "",
		UserName:  "",
		CreatedAt: "",
		UpdatedAt: "",
	}

	err := mockUser.EnsureCreationAttributes()
	meta := trace.Guarantee(err).Meta

	assert := test.Assert{T: t}
	assert.MapHasKey(meta, "fields")
	assert.IsString(meta["fields"])
	assert.MatchesString(meta["fields"].(string), "[Email UserName]")
}

func TestEnsureCreationAttributesPartialUser(t *testing.T) {
	mockUser := user.User{
		Id:        "",
		Email:     "email",
		UserName:  "",
		CreatedAt: "",
		UpdatedAt: "",
	}

	err := mockUser.EnsureCreationAttributes()
	meta := trace.Guarantee(err).Meta

	assert := test.Assert{T: t}
	assert.MapHasKey(meta, "fields")
	assert.IsString(meta["fields"])
	assert.MatchesString(meta["fields"].(string), "[UserName]")
}

func TestEnsureCreationAttributesCompleteUser(t *testing.T) {
	mockUser := user.User{
		Id:        "",
		Email:     "foo",
		UserName:  "foo",
		CreatedAt: "",
		UpdatedAt: "",
	}

	err := mockUser.EnsureCreationAttributes()

	assert := test.Assert{T: t}
	assert.ErrorIsNil(err)
}
