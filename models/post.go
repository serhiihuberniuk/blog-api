package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"time"
)

type Post struct {
	Id          string
	Title       string
	Description string
	CreatedBy   User
	CreatedAt   time.Time
	Tags        []string
	Comments    []Comment
}

func (p Post) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.Title, validation.Required, validation.Length(1, 50)),
		validation.Field(&p.Description, validation.Required))

}
