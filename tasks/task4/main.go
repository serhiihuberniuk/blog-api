package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/serhiihuberniuk/blog-api/tasks/task4/service/parser"
	"github.com/serhiihuberniuk/blog-api/tasks/task4/service/printer"
)

func main() {
	ctx := context.Background()
	client := http.Client{}
	p := parser.NewParser()
	pr := printer.NewPrinter()

	if err := printDailyNewsFromFootballSite(ctx, client, "https://football.ua/", p, pr); err != nil {
		log.Fatalf("cannot print news: %v", err)
	}

}

func printDailyNewsFromFootballSite(ctx context.Context, client http.Client, url string,
	parser *parser.Parser, printer *printer.Printer) error {
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("error while getting response: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response code is not 200 OK")
	}
	news, err := parser.GetDailyNewsFromSite(ctx, resp.Body)
	if err != nil {
		return fmt.Errorf("cannot get news from web site: %v", err)
	}

	for _, v := range news {
		printer.PrintNews(ctx, v)
	}

	return nil
}
