package main

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func ScrapeArticle(url string) (string, string, string, error) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", "", "", err
	}

	// biar dianggap browser asli
	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0 Safari/537.36",
	)

	resp, err := client.Do(req)
	if err != nil {
		return "", "", "", err
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", "", "", err
	}

	// ambil title
	title := strings.TrimSpace(doc.Find("title").First().Text())

	// ambil isi artikel
	var paragraphs []string

	doc.Find("p").Each(func(i int, s *goquery.Selection) {

		text := strings.TrimSpace(s.Text())

		// filter paragraf sampah
		if len(text) > 60 {
			paragraphs = append(paragraphs, text)
		}
	})

	content := strings.Join(paragraphs, "\n\n")

	// ambil og:image
	image := ""

	ogImage, exists := doc.Find(`meta[property="og:image"]`).Attr("content")

	if exists {
		image = ogImage
	}

	return title, content, image, nil
}
