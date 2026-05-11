package main

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ScrapeArticle(url string) (string, string, string, error) {

	resp, err := http.Get(url)
	if err != nil {
		return "", "", "", err
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", "", "", err
	}

	title := strings.TrimSpace(doc.Find("title").Text())

	var paragraphs []string

	doc.Find("p").Each(func(i int, s *goquery.Selection) {

		text := strings.TrimSpace(s.Text())

		if len(text) > 60 {
			paragraphs = append(paragraphs, text)
		}
	})

	content := strings.Join(paragraphs, "\n\n")

	image := ""

	ogImage, exists := doc.Find(`meta[property="og:image"]`).Attr("content")

	if exists {
		image = ogImage
	}

	return title, content, image, nil
}
