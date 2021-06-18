package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/pkg/errors"
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
		return errors.Wrap(errors.Cause(err), "validation failed")
	}

	return nil
}
