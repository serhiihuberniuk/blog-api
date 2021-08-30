package models

import (
	"testing"
)

func TestPost_Validate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		in         Post
		errMessage string
	}{
		{
			name: "Validation passed",
			in: Post{
				Title:       "title",
				Description: "description",
			},
			errMessage: "",
		},
		{
			name: "Title is empty",
			in: Post{
				Description: "description",
			},
			errMessage: requiredErrMessage,
		},
		{
			name: "Title is too long",
			in: Post{
				Title:       "titletitletitletitletitletitletitletitletitletitletitle",
				Description: "description",
			},
			errMessage: lengthErrMessage,
		},
		{
			name: "Description is empty",
			in: Post{
				Title: "title",
			},
			errMessage: requiredErrMessage,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.in.Validate()
			checkValidateErrorMessage(t, tc.errMessage, err)
		})
	}
}
