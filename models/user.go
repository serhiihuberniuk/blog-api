package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	//WARN [runner] Can't run linter goanalysis_metalinter: buildir: failed to load package :
	//could not load export data: no export data for "github.com/asaskevich/govalidator"
)

type User struct {
	ID        string
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (u User) Validate() error {
	return validation.ValidateStruct(&u,
		validation.Field(&u.Name, validation.Required, validation.Length(1, 30)),
		validation.Field(&u.Email, validation.Required, validation.Length(1, 30), is.Email),
	)

}
