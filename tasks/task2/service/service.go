package service

import (
	"context"
	"fmt"
)

type Service struct {
	storage storage
	buffer  buffer
}

func NewService(s storage, b buffer) *Service {
	return &Service{
		storage: s,
		buffer:  b,
	}
}

type buffer interface {
	BufferByCategory(category string, value int) []int
}

type storage interface {
	Save(ctx context.Context, category string, values []int) error
}

func (s *Service) Save(ctx context.Context, category string, value int) error {
	values := s.buffer.BufferByCategory(category, value)
	if len(values) > 0 {
		err := s.storage.Save(ctx, category, values)
		if err != nil {
			return fmt.Errorf("error while saving data to storage: %w", err)
		}
	}

	return nil
}
