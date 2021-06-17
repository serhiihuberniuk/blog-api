package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Post struct {
	ID          string
	Title       string
	Description string
	CreatedBy   string
	CreatedAt   time.Time
	Tags        []string
}

func (p Post) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Title, validation.Required, validation.Length(1, 50)),
		validation.Field(&p.Description, validation.Required))

}
