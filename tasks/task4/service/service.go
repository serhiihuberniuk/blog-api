package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
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
			return
		}
	})

	if err != nil {
		return fmt.Errorf("error occured while saving news to storage: %w", err)
	}

	return nil
}

func (s *Service) GetDailyNewsFromStorage(ctx context.Context) ([]models.FootballNews, error) {
	news, err := s.storage.GetAllNews(ctx)
	if err != nil {
		return nil, fmt.Errorf("error occures while get news from storage: %w", err)
	}

	return news, nil
}

func (s *Service) PrintDailyNews(_ context.Context, news []models.FootballNews) error {
	for _, footballNews := range news {
		fmt.Printf("Title: %s\nLink: %s\n\n", footballNews.Title, footballNews.Link)
	}

	return nil
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
