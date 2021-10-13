package parser

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/serhiihuberniuk/blog-api/tasks/task4/models"
)

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) GetDailyNewsFromSite(_ context.Context, r io.Reader) ([]models.FootballNews, error) {
	var footballNews models.FootballNews

	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, fmt.Errorf("error occurred while creating document: %w", err)
	}

	var allNews []models.FootballNews
	doc.Find(".news-feed.main-news").Find("ul").
		Find("li").Each(func(i int, selection *goquery.Selection) {
		link, ok := selection.Find("a").Attr("href")
		if !ok {
			return
		}
		footballNews.Link = link
		title := selection.Find("a").Text()
		footballNews.Title = formatTitle(title)

		allNews = append(allNews, footballNews)
	})

	return allNews, nil
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
