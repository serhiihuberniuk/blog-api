package buffer

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

type Buffer struct {
	bufferMutex   *sync.Mutex
	bufferStorage map[string]*categoryBuffer
	duration      time.Duration
	size          int
}

func NewBuffer(d time.Duration, size int) (*Buffer, error) {
	err := validateBufferParams(d, size)
	if err != nil {
		return nil, fmt.Errorf("error occured while creating buffer : %w", err)
	}
	return &Buffer{
		bufferMutex:   &sync.Mutex{},
		bufferStorage: make(map[string]*categoryBuffer),
		duration:      d,
		size:          size,
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

func (b *Buffer) BufferByCategory(category string, value int) []int {
	b.bufferMutex.Lock()

	cb, ok := b.bufferStorage[category]
	if ok {
		if b.put(cb, value) {
			b.bufferMutex.Unlock()

			return nil
		}

		b.bufferMutex.Unlock()
		return b.BufferByCategory(category, value)
	}

	cb = b.newCategoryBuffer(category)
	b.put(cb, value)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go cb.waitForFlush(wg)

	b.bufferMutex.Unlock()
	wg.Wait()

	return b.flushBuffer(category, cb)
}

type categoryBuffer struct {
	values []int
	timer  *time.Timer
	flush  chan struct{}
}

func (b *Buffer) put(cb *categoryBuffer, value int) bool {
	if len(cb.values) > (b.size - 1) {
		return false
	}

	cb.values = append(cb.values, value)

	if len(cb.values) == b.size {
		close(cb.flush)

		return true
	}

	if !cb.timer.Stop() {
		<-cb.timer.C
	}

	cb.timer.Reset(b.duration)

	return true
}

func (b *Buffer) newCategoryBuffer(category string) *categoryBuffer {
	b.bufferStorage[category] = &categoryBuffer{
		values: make([]int, 0, b.size),
		timer:  time.NewTimer(b.duration),
		flush:  make(chan struct{}),
	}

	return b.bufferStorage[category]
}

func (b *Buffer) flushBuffer(key string, cb *categoryBuffer) []int {
	values := cb.values

	b.bufferMutex.Lock()
	delete(b.bufferStorage, key)
	b.bufferMutex.Unlock()

	return values
}

func (cb *categoryBuffer) waitForFlush(wg *sync.WaitGroup) {
	select {
	case <-cb.timer.C:
		fmt.Println("time out")
	case <-cb.flush:
		fmt.Println("buffer full")
	}

	wg.Done()
}
