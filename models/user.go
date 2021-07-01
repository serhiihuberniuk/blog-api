package models

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

const maxLength = 30

type User struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateUserPayload struct {
	Name  string
	Email string
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
	)
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}
