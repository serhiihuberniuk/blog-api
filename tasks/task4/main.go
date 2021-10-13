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
	p := parser.NewParser(client)
	pr := printer.NewPrinter()

	if err := printDailyNewsFromFootballSite(ctx, "https://football.ua/", p, pr); err != nil {
		log.Fatalf("cannot print news: %v", err)
	}

}

func printDailyNewsFromFootballSite(ctx context.Context, url string, parser *parser.Parser, printer *printer.Printer) error {
	news, err := parser.GetDailyNewsFromSite(ctx, url)
	if err != nil {
		return fmt.Errorf("cannot get news from web site: %v", err)
	}

	for _, v := range news {
		printer.PrintNews(ctx, v)
	}

	return nil
}
