package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/serhiihuberniuk/blog-api/tasks/task2/buffer"
	"github.com/serhiihuberniuk/blog-api/tasks/task2/service"
	"github.com/serhiihuberniuk/blog-api/tasks/task2/storage"
)

const (
	bufferDuration = time.Second * 3
	bufferSize     = 5

	category = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	n        = 1000
)

func main() {
	ctx := context.Background()

	st := storage.NewStorage()
	b, err := buffer.NewBuffer(st, bufferDuration, bufferSize)
	if err != nil {
		log.Fatal(err)
	}
	s := service.NewService(b)

	errs := make(chan error, 1)
	go func() {
		for err := range errs {
			fmt.Println(err)
		}
	}()

	for i := 0; i < n; i++ {
		go func() {
			if err := s.Save(ctx, string(category[rand.Intn(len(category))]), rand.Intn(n)); err != nil {
				errs <- err
			}
		}()
	}

	b.Wait()
	fmt.Println(st.Count())
}
