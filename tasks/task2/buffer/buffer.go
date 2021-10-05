package buffer

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type storage interface {
	Save(ctx context.Context, category string, values []int) error
}

type Buffer struct {
	storage         storage
	lock            *sync.RWMutex
	numberOfBuffers int32
	buffers         map[string]chan int

	duration time.Duration
	size     int
}

func NewBuffer(storage storage, d time.Duration, size int) (*Buffer, error) {
	err := validateBufferParams(d, size)
	if err != nil {
		return nil, fmt.Errorf("error occured while creating buffer : %w", err)
	}

	return &Buffer{
		storage:  storage,
		lock:     &sync.RWMutex{},
		buffers:  make(map[string]chan int),
		duration: d,
		size:     size,
	}, nil
}

func validateBufferParams(d time.Duration, size int) error {
	if d < time.Second {
		return errors.New("validation error: buffer duration is too short")
	}

	if size < 2 {
		return errors.New("validation error: buffer is too short")
	}

	return nil
}

func (b *Buffer) BufferByCategory(_ context.Context, category string, value int) error {

	b.lock.Lock()
	defer b.lock.Unlock()

	ch, ok := b.buffers[category]
	if !ok {
		ch = make(chan int)
		b.buffers[category] = ch
		b.newCategory(category, value, ch)

		return nil
	}

	ch <- value

	return nil

}

func (b *Buffer) Wait() {
	for {
		if atomic.LoadInt32(&b.numberOfBuffers) == 0 {
			return
		}
	}
}

func (b *Buffer) newCategory(category string, value int, ch chan int) {

	go func(ch chan int) {
		values := make([]int, 0, b.size)
		values = append(values, value)
		atomic.AddInt32(&b.numberOfBuffers, 1)

		flush := func() {
			if len(values) == 0 {
				return
			}

			if err := b.storage.Save(context.Background(), category, values); err != nil {
				fmt.Println("Error during saving buffer: ", err.Error())
			}
			atomic.AddInt32(&b.numberOfBuffers, -1*int32(len(values)))

			values = values[0:0]

		}

		for {
			select {
			case v := <-ch:
				values = append(values, v)
				atomic.AddInt32(&b.numberOfBuffers, 1)
				if len(values) == b.size {
					flush()
					fmt.Println("full")

				}
			case <-time.After(b.duration):
				flush()
				fmt.Println("timer")

			}
		}
	}(ch)

}
