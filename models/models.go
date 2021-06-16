package models

import "time"

//Моделі:
//	User: id, name, email, created_at, updated_at
//	Post: id, title, description, created_by, created_at, tags
//	Comment: id, content, created_by, created_at

type User struct {
	Id        string
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Post struct {
	Id          string
	Title       string
	Description string
	CreatedBy   User
	CreatedAt   time.Time
	Tags        []string
}

type Comment struct {
	Id        string
	Content   string
	CreatedBy User
	CreatedAt time.Time
}
