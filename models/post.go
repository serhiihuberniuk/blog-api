package models

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

const maxLengthTitle = 50

const (
	FilterPostsByTitle     FilterPostsByField = "title"
	FilterPostsByCreatedBy FilterPostsByField = "created_by"
	FilterPostsByTags      FilterPostsByField = "tag"
)

const (
	SortPostsByTitle     SortPostsByField = "title"
	SortPostsByTags      SortPostsByField = "tag"
	SortPostsByCreatedAt SortPostsByField = "created_at"
)

type Post struct {
	ID          string    `bson:"_id,omitempty"`
	Title       string    `bson:"title,omitempty"`
	Description string    `bson:"description,omitempty"`
	CreatedBy   string    `bson:"created_by,omitempty"`
	CreatedAt   time.Time `bson:"created_at,omitempty"`
	Tags        []string  `bson:"tags,omitempty"`
}

type CreatePostPayload struct {
	Title       string
	Description string
	Tags        []string
	AuthorID    string
}

type UpdatePostPayload struct {
	PostID      string
	Title       string
	Description string
	Tags        []string
}

type FilterPostsByField string

type FilterPosts struct {
	Field FilterPostsByField
	Value string
}

type SortPostsByField string

type SortPosts struct {
	SortByField SortPostsByField
	IsASC       bool
}

func (p *Post) Validate() error {
	err := validation.ValidateStruct(p,
		validation.Field(&p.Title, validation.Required, validation.Length(1, maxLengthTitle)),
		validation.Field(&p.Description, validation.Required))
	if err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	return nil
}
