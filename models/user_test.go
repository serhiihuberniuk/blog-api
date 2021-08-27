package models

import (
	"testing"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	testCases := []struct {
		name    string
		in      User
		wantErr bool
	}{
		{
			name: "Validation passed",
			in: User{
				Name:     "name",
				Email:    "email@mail.com",
				Password: "password",
			},
			wantErr: false,
		},
		{
			name: "Name field is empty",
			in: User{
				Email:    "email@mail.com",
				Password: "password",
			},
			wantErr: true,
		},
		{
			name: "Name field is too short",
			in: User{
				Name:     "",
				Email:    "email@mail.com",
				Password: "password",
			},
			wantErr: true,
		},
		{
			name: "Name field is too long",
			in: User{
				Name:     "ssssssssssssssssssssssssssssssssssssssssss",
				Email:    "email@mail.com",
				Password: "password",
			},
			wantErr: true,
		},
		{
			name: "Email field is empty",
			in: User{
				Name:     "name",
				Password: "password",
			},
			wantErr: true,
		},
		{
			name: "Email field is too short",
			in: User{
				Name:     "name",
				Email:    "",
				Password: "password",
			},
			wantErr: true,
		},
		{
			name: "Email field is too long",
			in: User{
				Name:     "",
				Email:    "sssssssssssss@ssssssssssssssmail.sssssssssssssss",
				Password: "password",
			},
			wantErr: true,
		},
		{
			name: "Email field is not an email",
			in: User{
				Name:     "name",
				Email:    "not an email address",
				Password: "password",
			},
			wantErr: true,
		},
		{
			name: "Password field is empty",
			in: User{
				Name:  "name",
				Email: "email@mail.com",
			},
			wantErr: true,
		},
		{
			name: "Password is too short",
			in: User{
				Name:     "name",
				Email:    "",
				Password: "short",
			},
			wantErr: true,
		},
		{
			name: "Password is too long",
			in: User{
				Name:     "name",
				Email:    "email@mail.com",
				Password: "passwordpasswordpasswordpassword",
			},
			wantErr: true,
		},
	}

	for _, tc := range testCases {
		var isError = false
		t.Run(tc.name, func(t *testing.T) {
			err := validation.ValidateStruct(&tc.in,
				validation.Field(&tc.in.Name, validation.Required, validation.Length(1, maxLength)),
				validation.Field(&tc.in.Email, validation.Required, validation.Length(1, maxLength), is.Email),
				validation.Field(&tc.in.Password, validation.Required, validation.Length(minLengthOfPassword, maxLength)),
			)
			if err != nil {
				isError = true
			}
			assert.Equal(t, tc.wantErr, isError)

		})
	}
}
