package models

import (
	"testing"
)

func TestComment_Validate(t *testing.T) {
	testCases := []struct {
		name       string
		in         Comment
		errMessage string
	}{
		{
			name: "Validation passed",
			in: Comment{
				Content: "content",
			},
			errMessage: "",
		},
		{
			name:       "Content is empty",
			in:         Comment{},
			errMessage: requiredErrMessage,
		},
	}

	for _, tc := range testCases {

		t.Run(tc.name, func(t *testing.T) {
			err := tc.in.Validate()
			checkValidateErrorMessage(t, tc.errMessage, err)
		})
	}
}
