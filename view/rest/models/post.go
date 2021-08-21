package models

import (
	"time"
)

type CreatePostRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}

type GetPostResponse struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatedBy   string    `json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt"`
	Tags        []string  `json:"tags"`
}

type UpdatePostRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
}
