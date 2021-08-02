package models

import "errors"

var (
	UserNotFound    = ErrorNotFound{Item: "User"}
	PostNotFound    = ErrorNotFound{Item: "Post"}
	CommentNotFound = ErrorNotFound{Item: "Comment"}
	ErrorBadRequest = errors.New("request data is invalid")
)

type ErrorNotFound struct {
	Item string
}

func (e ErrorNotFound) Error() string {
	return e.Item + " is not found with such ID"
}
