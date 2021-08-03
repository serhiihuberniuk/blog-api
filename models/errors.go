package models

import "errors"

var (
	ErrNotFoundUser    = errors.New("user is not found with such ID")
	ErrNotFoundPost    = errors.New("user is not found with such ID")
	ErrNotFoundComment = errors.New("user is not found with such ID")
)
