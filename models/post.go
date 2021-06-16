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
	return validation.ValidateStruct(&p)//todo

}
