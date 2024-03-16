package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SetSubject(t *testing.T) {
	user := User{}
	user.ID = "asdf"
	origin := ImportOrigin{}
	if assert.NoError(t, origin.SetSubject(user)) {
		assert.Equal(t, "User", origin.SubjectType)
		assert.Equal(t, "asdf", origin.SubjectID)
	}
}

func Test_SetSubject_SubjectIsNoModel(t *testing.T) {
	type NotABaseModel struct {
		Identifier string
	}
	user := NotABaseModel{
		Identifier: "asdf",
	}
	origin := ImportOrigin{}
	assert.EqualError(t, origin.SetSubject(user), "subject of type model.NotABaseModel does not have an ID field")
}
