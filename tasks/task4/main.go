package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/serhiihuberniuk/blog-api/tasks/task4/service"
	"github.com/serhiihuberniuk/blog-api/tasks/task4/storage"
)

func main() {
	client := http.Client{}
	st := storage.NewStorage()
	s := service.NewService(st)

	resp, err := getResponse(client, "https://football.ua/")
	if err != nil {
		log.Fatal("cannot get response")
	}
	defer resp.Body.Close()

	err = s.GetDailyNews(resp.Body)
	if err != nil {
		log.Fatalf("cannot get news: %v", err)
	}

	s.PrintDailyNews()
}

func getResponse(client http.Client, url string) (*http.Response, error) {
	resp, err := client.Get(url)
	if err != nil {
		return nil, fmt.Errorf("error while getting response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("response code is not 200 OK")
	}

	return resp, nil
}
