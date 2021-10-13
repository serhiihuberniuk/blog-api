package storage

import (
	"context"

	"github.com/serhiihuberniuk/blog-api/tasks/task4/models"
)

type Storage struct {
	newsStorage []models.FootballNews
}

func NewStorage() *Storage {
	return &Storage{
		newsStorage: []models.FootballNews{},
	}
}

func (s *Storage) SaveNew(_ context.Context, footballNews models.FootballNews) error {
	s.newsStorage = append(s.newsStorage, footballNews)

	return nil
}

func (s *Storage) GetAllNews(_ context.Context) ([]models.FootballNews, error) {
	return s.newsStorage, nil
}
