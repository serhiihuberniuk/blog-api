package models

import (
	"testing"
)

func TestUser_Validate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		in         User
		errMessage string
	}{
		{
			name: "Validation passed",
			in: User{
				Name:     "name",
				Email:    "email@mail.com",
				Password: "password",
			},
			errMessage: "",
		},
		{
			name: "Name field is empty",
			in: User{
				Email:    "email@mail.com",
				Password: "password",
			},
			errMessage: requiredErrMessage,
		},
		{
			name: "Name field is too short",
			in: User{
				Name:     "n",
				Email:    "email@mail.com",
				Password: "password",
			},
			errMessage: lengthErrMessage,
		},
		{
			name: "Name field is too long",
			in: User{
				Name:     "ssssssssssssssssssssssssssssssssssssssssss",
				Email:    "email@mail.com",
				Password: "password",
			},
			errMessage: lengthErrMessage,
		},
		{
			name: "Email field is empty",
			in: User{
				Name:     "name",
				Password: "password",
			},
			errMessage: requiredErrMessage,
		},
		{
			name: "Email field is too short",
			in: User{
				Name:     "name",
				Email:    "m@m.com",
				Password: "password",
			},
			errMessage: lengthErrMessage,
		},
		{
			name: "Email field is too long",
			in: User{
				Name:     "name",
				Email:    "sssssssssssss@ssssssssssssssmail.sssssssssssssss",
				Password: "password",
			},
			errMessage: lengthErrMessage,
		},
		{
			name: "Email field is not an email",
			in: User{
				Name:     "name",
				Email:    "not an email address",
				Password: "password",
			},
			errMessage: isMailErrMessage,
		},
		{
			name: "Password field is empty",
			in: User{
				Name:  "name",
				Email: "email@mail.com",
			},
			errMessage: requiredErrMessage,
		},
		{
			name: "Password is too short",
			in: User{
				Name:     "name",
				Email:    "email@mail.com",
				Password: "short",
			},
			errMessage: lengthErrMessage,
		},
		{
			name: "Password is too long",
			in: User{
				Name:     "name",
				Email:    "email@mail.com",
				Password: "passwordpasswordpasswordpassword",
			},
			errMessage: lengthErrMessage,
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
