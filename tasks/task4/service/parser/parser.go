package parser

import (
	"bytes"
	"context"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/serhiihuberniuk/blog-api/tasks/task4/models"
)

type Parser struct {
	source source
}

func NewParser(s source) *Parser {
	return &Parser{
		source: s,
	}
}

type source interface {
	GetDOMContent(url string) ([]byte, error)
}

func (p *Parser) ParseDailyNews(_ context.Context, url string) ([]models.FootballNews, error) {
	var footballNews models.FootballNews

	content, err := p.source.GetDOMContent(url)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(content))
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
