package main

import (
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type ScrapeResult struct {
	FinalURL string `json:"final_url"`
	Title    string `json:"title"`
	Image    string `json:"image"`
	Content  string `json:"content"`
}

// ScrapeArticle
func ScrapeArticle(url string) (*ScrapeResult, error) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return nil, err
	}

	// User-Agent browser asli
	req.Header.Set(
		"User-Agent",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/136.0 Safari/537.36",
	)

	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Accept", "text/html,application/xhtml+xml")

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	// URL final setelah redirect
	finalURL := resp.Request.URL.String()

	// baca seluruh html
	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	htmlContent := string(bodyBytes)

	// parse html
	doc, err := goquery.NewDocumentFromReader(
		strings.NewReader(htmlContent),
	)

	if err != nil {
		return nil, err
	}

	// DEBUG DIV
	doc.Find("div").Each(func(i int, s *goquery.Selection) {

		class, _ := s.Attr("class")

		text := strings.TrimSpace(s.Text())

		if len(text) > 500 {

			println("CLASS:", class)

			if len(text) > 300 {
				println("TEXT:", text[:300])
			} else {
				println("TEXT:", text)
			}

			println("================================")
		}
	})

	// =========================
	// TITLE
	// =========================

	title := strings.TrimSpace(
		doc.Find("title").First().Text(),
	)

	// =========================
	// CONTENT
	// =========================

	var paragraphs []string

	selectors := []string{
		"article p",
		".read__content p",
		".content-text p",
		".article-content p",
		".detail__body-text p",
		".post-content p",
		".entry-content p",
		".content p",
		"p",
	}

	found := false

	for _, selector := range selectors {

		doc.Find(selector).Each(func(i int, s *goquery.Selection) {

			text := strings.TrimSpace(s.Text())

			// filter paragraf pendek
			if len(text) > 60 {

				paragraphs = append(paragraphs, text)

				found = true
			}
		})

		if found {
			break
		}
	}

	content := strings.Join(paragraphs, "\n\n")

	// =========================
	// IMAGE
	// =========================

	image := ""

	ogImage, exists := doc.Find(
		`meta[property="og:image"]`,
	).Attr("content")

	if exists {
		image = ogImage
	}

	// =========================
	// DEBUG FALLBACK
	// =========================

	if content == "" {

		debugHTML := htmlContent

		if len(debugHTML) > 3000 {
			debugHTML = debugHTML[:3000]
		}

		content = debugHTML
	}

	return &ScrapeResult{
		FinalURL: finalURL,
		Title:    title,
		Image:    image,
		Content:  content,
	}, nil
}
