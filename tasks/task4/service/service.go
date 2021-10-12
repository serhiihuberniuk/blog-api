package service

import (
	"bytes"
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
	SaveNew(footballNew models.FootballNew)
	GetAllNews() []models.FootballNew
}

func (s *Service) GetDailyNews(r io.Reader) error {
	var footballNew models.FootballNew

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
		footballNew.Link = link
		title := selection.Find("a").Text()
		footballNew.Title = formatTitle(title)
		s.storage.SaveNew(footballNew)
	})

	return nil
}

func (s *Service) PrintDailyNews() {

	for _, footballNew := range s.storage.GetAllNews() {
		fmt.Printf("Title: %s\nLink: %s\n\n", footballNew.Title, footballNew.Link)
	}
}

func formatTitle(title string) string {
	var previousLetter bool
	var byteBuffer bytes.Buffer

	for _, letter := range title {
		currentLetter := letter == ' '
		if currentLetter {
			if !previousLetter {
				byteBuffer.WriteRune(letter)
			}
		} else {
			byteBuffer.WriteRune(letter)
		}

		previousLetter = currentLetter
	}

	return strings.Trim(byteBuffer.String(), " \n")
}
