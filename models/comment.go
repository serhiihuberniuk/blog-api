package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"time"
)

type Comment struct {
	Id        string
	Content   string
	CreatedBy User
	CreatedAt time.Time
	CreatedTo Post
}

func (c Comment) Validate() error {
	return validation.ValidateStruct(&c)//todo

}
