package models

import "errors"

var (
	ErrNotFound         = errors.New("not found")
	ErrNotAuthenticated = errors.New("not authenticated")
)
