package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type Comment struct {
	ID        string
	Content   string
	CreatedBy string
	CreatedAt time.Time
	PostID    string
}

func (c Comment) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.Content, validation.Required))

}
