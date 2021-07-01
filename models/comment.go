package models

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	FilterCommentsByAuthor    FilterCommentsByField = "CreatedBy"
	FilterCommentsByCreatedAt FilterCommentsByField = "CreatedAt"
	FilterCommentsByPost      FilterCommentsByField = "PostID"
)

const (
	SortCommentByCreatedAt SortCommentsByField = "CreatedAt"
)

type Comment struct {
	ID        string
	Content   string
	CreatedBy string
	CreatedAt time.Time
	PostID    string
}

type CreateCommentPayload struct {
	Content  string
	PostID   string
	AuthorID string
}

type UpdateCommentPayload struct {
	CommentID string
	Content   string
}

type FilterCommentsByField string

type FilterComments struct {
	Field FilterCommentsByField
	Value string
}

type SortCommentsByField string

type SortComments struct {
	Field SortCommentsByField
	IsASC bool
}

func (c *Comment) Validate() error {
	err := validation.ValidateStruct(c, validation.Field(&c.Content, validation.Required))
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}
