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

func (u *User) Validate() error {
	err := validation.ValidateStruct(u,
		validation.Field(&u.Name, validation.Required, validation.Length(1, maxLength)),
		validation.Field(&u.Email, validation.Required, validation.Length(1, maxLength), is.Email),
	)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
