package data

import (
	"testing"
)

func TestTagValidationFail(t *testing.T) {
	tag, _ := NewTag("some tag with s  pace")

	v := NewValidation()

	err := v.Validate(tag)

	if err == nil {

		t.Logf("Expected error wrong Format")
		t.FailNow()

	}

	t.Logf("Error %v", err.Errors())
	t.FailNow()

}

func TestTagValidationSuccess(t *testing.T) {
	tag, _ := NewTag("a_tag_without_white_space")

	v := NewValidation()

	err := v.Validate(tag)

	if err != nil {

		t.Logf("Expected error wrong Format")
		t.FailNow()

	}

}
