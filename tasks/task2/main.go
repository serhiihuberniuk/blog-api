//Припустимо в тебе є сервіс для зберігання якихось даних, він має один метод:
//- Save(ctx context.Context, category string, value int) error
//    - category = A, B, C ...всі великі літери
//    - value    = any int value
//
//Треба написати:
//    - Клієнт сервісу: клієнт паралельно запускає N паралельних процесів збереження даних в сервіс
//    - Сам сервіс: сервіс приймає дані методом Save і відправляєш їх в сторедж, але з наступною логікою:
//        - не просто відправляє, а робить буферезацію по категорії
//        - максимальний розмір буфера = 5 елементів
//        - максимальний час буфера = 3 секунди
//    - Сторедж: просто поки можна писати в консоль дані які сторедж має зберегти
//
//Створи окрему папку в проекті і напиши в ній дане рішення, використовуючи усі правила SOLID you know.
package main

import (
	"context"
	"log"
	"math/rand"
	"sync"
	"time"

	"github.com/serhiihuberniuk/blog-api/tasks/task2/buffer"
	"github.com/serhiihuberniuk/blog-api/tasks/task2/service"
	"github.com/serhiihuberniuk/blog-api/tasks/task2/storage"
)

const (
	bufferDuration = time.Second * 3
	bufferSize     = 10
)

var category = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func main() {
	ctx := context.Background()

	st := storage.NewStorage()
	b, err := buffer.NewBuffer(bufferDuration, bufferSize)
	if err != nil {
		log.Fatal(err)
	}
	s := service.NewService(st, b)

	var wg sync.WaitGroup

	for i := 0; i < 1001; i++ {
		randomIndex := rand.Intn(len(category))
		randomValue := rand.Intn(1000)

		wg.Add(1)

		go func() {
			err := s.Save(ctx, string(category[randomIndex]), randomValue)
			if err != nil {
				wg.Done()
				log.Fatal(err)
			}

			wg.Done()
		}()

	}

	wg.Wait()
}
