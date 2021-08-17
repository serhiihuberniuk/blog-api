package models

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

const (
	maxLength           = 30
	minLengthOfPassword = 8
)

type User struct {
	ID        string    `bson:"_id,omitempty"`
	Name      string    `bson:"name,omitempty"`
	Email     string    `bson:"email,omitempty"`
	CreatedAt time.Time `bson:"created_at,omitempty"`
	UpdatedAt time.Time `bson:"updated_at,omitempty"`
	Password  string    `bson:"password, omitempty"`
}

type CreateUserPayload struct {
	Name     string
	Email    string
	Password string
}

type UpdateUserPayload struct {
	UserID string
	Name   string
	Email  string
}

func (u *User) Validate() error {
	err := validation.ValidateStruct(u,
		validation.Field(&u.Name, validation.Required, validation.Length(1, maxLength)),
		validation.Field(&u.Email, validation.Required, validation.Length(1, maxLength), is.Email),
		validation.Field(&u.Password, validation.Required, validation.Length(minLengthOfPassword, maxLength)),
	)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}
