package service

import (
	"context"
)

type Service struct {
	buffer buffer
}

type buffer interface {
	BufferByCategory(ctx context.Context, category string, value int) error
}

func NewService(b buffer) *Service {
	return &Service{
		buffer: b,
	}
}

func (s *Service) Save(ctx context.Context, category string, value int) error {
	return s.buffer.BufferByCategory(ctx, category, value)
}
