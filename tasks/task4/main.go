package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	err := getDailyNews()
	if err != nil {
		log.Fatal(err)
	}
}

func getDailyNews() error {
	client := http.Client{}

	resp, err := client.Get("https://football.ua/")
	if err != nil {
		return fmt.Errorf("error while getting response: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New("response code is not 200 OK")
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return fmt.Errorf("error while creating document from response: %w", err)
	}

	doc.Find(".news-feed.main-news").Find("ul").
		Find("li").Each(func(i int, selection *goquery.Selection) {
		link, _ := selection.Find("a").Attr("href")
		text := selection.Find("a").Text()
		text = formatText(text)
		fmt.Printf("Title: %s\nLink: %s\n\n", text, link)
	})

	return nil
}

func formatText(text string) string {
	var previousLetter bool
	var byteBuffer bytes.Buffer
	for _, letter := range text {
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
