package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_ModelAndID_NotAStruct(t *testing.T) {
	value := "asd"
	model, id, err := ModelAndID(value)
	assert.Equal(t, "", model)
	assert.Equal(t, "", id)
	assert.ErrorContains(t, err, "Value of type string is not a struct")
}

type Peter struct {
	Foo string
	Bar string
}

func Test_ModelAndID_StructWithoutID(t *testing.T) {
	value := Peter{Foo: "FOO", Bar: "BAR"}
	model, id, err := ModelAndID(value)
	assert.Equal(t, "Peter", model)
	assert.Equal(t, "", id)
	assert.ErrorContains(t, err, "Value of type model.Peter does not contain field \"ID\"")
}

type Joanna struct {
	ID  string
	Foo string
	Bar string
}

func Test_ModelAndID_StructWithID(t *testing.T) {
	value := Joanna{ID: "c730e06d-4ece-47ee-927a-1961b6c3367e", Foo: "FOO", Bar: "BAR"}
	model, id, err := ModelAndID(value)
	assert.Equal(t, "Joanna", model)
	assert.Equal(t, "c730e06d-4ece-47ee-927a-1961b6c3367e", id)
	assert.NoError(t, err)
}
