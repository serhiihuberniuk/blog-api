package models

import (
	"fmt"
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

func (c *Comment) Validate() error {
	err := validation.ValidateStruct(c, validation.Field(&c.Content, validation.Required))
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	return nil
}
