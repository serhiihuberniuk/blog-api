package storage

import "github.com/serhiihuberniuk/blog-api/tasks/task4/models"

type Storage struct {
	newsStorage []models.FootballNew
}

func NewStorage() *Storage {
	newsStorage := make([]models.FootballNew, 0, 11)

	return &Storage{
		newsStorage: newsStorage,
	}
}

func (s *Storage) SaveNew(footballNew models.FootballNew) {
	s.newsStorage = append(s.newsStorage, footballNew)
}

func (s *Storage) GetAllNews() []models.FootballNew {
	return s.newsStorage
}
