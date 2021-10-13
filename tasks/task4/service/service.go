package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/serhiihuberniuk/blog-api/tasks/task4/models"
)

type Service struct {
	storage storage
}

func NewService(s storage) *Service {
	return &Service{
		storage: s,
	}
}

type storage interface {
	SaveNew(ctx context.Context, footballNews models.FootballNews) error
	GetAllNews(ctx context.Context) ([]models.FootballNews, error)
}

func (s *Service) SaveDailyNews(ctx context.Context, r io.Reader) error {
	var footballNews models.FootballNews

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return fmt.Errorf("error occurred while creating document: %w", err)
	}

	doc.Find(".news-feed.main-news").Find("ul").
		Find("li").Each(func(i int, selection *goquery.Selection) {
		link, ok := selection.Find("a").Attr("href")
		if !ok {
			return
		}
		footballNews.Link = link
		title := selection.Find("a").Text()
		footballNews.Title = formatTitle(title)

		err = s.storage.SaveNew(ctx, footballNews)
		if err != nil {
			log.Println("error occurred while saving news:", err)
		}
	})

	return nil
}

func (s *Service) PrintNews(news models.FootballNews) {
	fmt.Printf("Title: %s\nLink: %s\n\n", news.Title, news.Link)
}

func formatTitle(title string) string {
	var isPreviousLetterGap bool
	var byteBuffer bytes.Buffer

	for _, letter := range title {
		isCurrentLetterGap := letter == ' '
		if isCurrentLetterGap {
			if !isPreviousLetterGap {
				byteBuffer.WriteRune(letter)
			}
		} else {
			byteBuffer.WriteRune(letter)
		}

		isPreviousLetterGap = isCurrentLetterGap
	}

	return strings.Trim(byteBuffer.String(), " \n")
}
