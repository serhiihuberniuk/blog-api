// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

import (
	"fmt"
	"io"
	"strconv"
)

type Comment struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	CreatedBy *User  `json:"createdBy"`
	AuthorID  string `json:"authorID"`
	CreatedAt string `json:"createdAt"`
	Post      *Post  `json:"post"`
	PostID    string `json:"postID"`
}

type CreateCommentInput struct {
	Content   string `json:"content"`
	CreatedBy string `json:"createdBy"`
	PostID    string `json:"postId"`
}

type CreatePostInput struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	CreatedBy   string   `json:"createdBy"`
}

type CreateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type FilterCommentsInput struct {
	Field FilterCommentsField `json:"field"`
	Value string              `json:"value"`
}

type FilterPostInput struct {
	Field FilterPostsField `json:"field"`
	Value string           `json:"value"`
}

type PaginationInput struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
}

type Post struct {
	ID          string   `json:"id"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	CreatedBy   *User    `json:"createdBy"`
	AuthorID    string   `json:"authorID"`
	CreatedAt   string   `json:"createdAt"`
	Tags        []string `json:"tags"`
}

type SortCommentsInput struct {
	Field SortCommentsField `json:"field"`
	IsAsc bool              `json:"isAsc"`
}

type SortPostsInput struct {
	Field SortPostsField `json:"field"`
	IsAsc bool           `json:"isAsc"`
}

type UpdateCommentInput struct {
	Content string `json:"content"`
}

type UpdatePostInput struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type UpdateUserInput struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type User struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

type FilterCommentsField string

const (
	FilterCommentsFieldCreatedBy FilterCommentsField = "CREATED_BY"
	FilterCommentsFieldPostID    FilterCommentsField = "POST_ID"
	FilterCommentsFieldCreatedAt FilterCommentsField = "CREATED_AT"
)

var AllFilterCommentsField = []FilterCommentsField{
	FilterCommentsFieldCreatedBy,
	FilterCommentsFieldPostID,
	FilterCommentsFieldCreatedAt,
}

func (e FilterCommentsField) IsValid() bool {
	switch e {
	case FilterCommentsFieldCreatedBy, FilterCommentsFieldPostID, FilterCommentsFieldCreatedAt:
		return true
	}
	return false
}

func (e FilterCommentsField) String() string {
	return string(e)
}

func (e *FilterCommentsField) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = FilterCommentsField(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid FilterCommentsField", str)
	}
	return nil
}

func (e FilterCommentsField) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type FilterPostsField string

const (
	FilterPostsFieldCreatedBy FilterPostsField = "CREATED_BY"
	FilterPostsFieldTitle     FilterPostsField = "TITLE"
	FilterPostsFieldTags      FilterPostsField = "TAGS"
)

var AllFilterPostsField = []FilterPostsField{
	FilterPostsFieldCreatedBy,
	FilterPostsFieldTitle,
	FilterPostsFieldTags,
}

func (e FilterPostsField) IsValid() bool {
	switch e {
	case FilterPostsFieldCreatedBy, FilterPostsFieldTitle, FilterPostsFieldTags:
		return true
	}
	return false
}

func (e FilterPostsField) String() string {
	return string(e)
}

func (e *FilterPostsField) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = FilterPostsField(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid FilterPostsField", str)
	}
	return nil
}

func (e FilterPostsField) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type SortCommentsField string

const (
	SortCommentsFieldCreatedAt SortCommentsField = "CREATED_AT"
)

var AllSortCommentsField = []SortCommentsField{
	SortCommentsFieldCreatedAt,
}

func (e SortCommentsField) IsValid() bool {
	switch e {
	case SortCommentsFieldCreatedAt:
		return true
	}
	return false
}

func (e SortCommentsField) String() string {
	return string(e)
}

func (e *SortCommentsField) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SortCommentsField(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SortCommentsField", str)
	}
	return nil
}

func (e SortCommentsField) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}

type SortPostsField string

const (
	SortPostsFieldCreatedAt SortPostsField = "CREATED_AT"
	SortPostsFieldTitle     SortPostsField = "TITLE"
)

var AllSortPostsField = []SortPostsField{
	SortPostsFieldCreatedAt,
	SortPostsFieldTitle,
}

func (e SortPostsField) IsValid() bool {
	switch e {
	case SortPostsFieldCreatedAt, SortPostsFieldTitle:
		return true
	}
	return false
}

func (e SortPostsField) String() string {
	return string(e)
}

func (e *SortPostsField) UnmarshalGQL(v interface{}) error {
	str, ok := v.(string)
	if !ok {
		return fmt.Errorf("enums must be strings")
	}

	*e = SortPostsField(str)
	if !e.IsValid() {
		return fmt.Errorf("%s is not a valid SortPostsField", str)
	}
	return nil
}

func (e SortPostsField) MarshalGQL(w io.Writer) {
	fmt.Fprint(w, strconv.Quote(e.String()))
}
