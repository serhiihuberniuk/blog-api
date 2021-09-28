package storage

import (
	"context"
	"fmt"
)

var i int

type Storage struct{}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) Save(ctx context.Context, category string, values []int) error {
	fmt.Println(category, values)
	i += len(values)
	fmt.Println(i)

	return nil
}
