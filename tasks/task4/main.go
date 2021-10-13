package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/serhiihuberniuk/blog-api/tasks/task4/service"
	"github.com/serhiihuberniuk/blog-api/tasks/task4/storage"
)

func main() {
	ctx := context.Background()
	client := http.Client{}
	st := storage.NewStorage()
	s := service.NewService(st)

	if err := printDailyNewsFromFootballSite(ctx, "https://football.ua/", client, s, st); err != nil {
		log.Fatalf("cannot print news: %v", err)
	}

}

func printDailyNewsFromFootballSite(ctx context.Context, url string, client http.Client,
	service *service.Service, storage *storage.Storage) error {
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("error while getting response: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response code is not 200 OK")
	}

	err = service.SaveDailyNews(ctx, resp.Body)
	if err != nil {
		return fmt.Errorf("cannot get news: %v", err)
	}

	news, err := storage.GetAllNews(ctx)
	if err != nil {
		return fmt.Errorf("cannot get news from storage: %v", err)
	}

	for _, v := range news {
		service.PrintNews(v)
	}

	return nil
}
