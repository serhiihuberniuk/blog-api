package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"time"
)

type User struct {
	Id        string
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
	Posts     []Post
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u)//todo

}
