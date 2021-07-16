package models

import (
	"time"
)

type CreateCommentRequest struct {
	Content  string `json:"content"`
	PostID   string `json:"postId"`
	AuthorID string `json:"authorId"`
}

type GetCommentResponse struct {
	ID        string    `json:"id"`
	Content   string    `json:"content"`
	CreatedBy string    `json:"createdBy"`
	CreatedAt time.Time `json:"createdAt"`
	PostID    string    `json:"postId"`
}

type UpdateCommentRequest struct {
	Content string `json:"content"`
}
