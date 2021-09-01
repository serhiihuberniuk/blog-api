package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	requiredErrMessage = "cannot be blank"
	lengthErrMessage   = "the length must be between"
	isMailErrMessage   = "must be a valid email address"
)

func checkValidateErrorMessage(t *testing.T, errMessage string, err error) {
	t.Helper()

	if errMessage == "" {
		assert.NoError(t, err)

		return
	}

	assert.Contains(t, err.Error(), errMessage)
}
