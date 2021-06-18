package models

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

const maxLengthTitle = 50

type Post struct {
	ID          string
	Title       string
	Description string
	CreatedBy   string
	CreatedAt   time.Time
	Tags        []string
}

func (p *Post) Validate() error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.Title, validation.Required, validation.Length(1, maxLengthTitle)),
		validation.Field(&p.Description, validation.Required))
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
