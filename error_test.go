package dry

import (
	"errors"
	"testing"
)

const TestErrorString = "TestError"

func Test_Error(t *testing.T) {
	err := AsError(TestErrorString)
	if err == nil || err.Error() != TestErrorString {
		t.Fail()
	}

	err = AsError(errors.New(TestErrorString))
	if err == nil || err.Error() != TestErrorString {
		t.Fail()
	}
}
