package storage

import (
	"context"
	"fmt"
	"sync/atomic"
)

type Storage struct {
	i int32
}

func NewStorage() *Storage {
	return &Storage{}
}

func (s *Storage) Save(_ context.Context, category string, values []int) error {
	fmt.Println(category, values)
	atomic.AddInt32(&s.i, int32(len(values)))

	return nil
}

func (s *Storage) Count() int32 {
	return atomic.LoadInt32(&s.i)
}
