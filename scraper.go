package main

import (
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// ScrapeArticle
func ScrapeArticle(url string) (string, string, string, error) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		return "", "", "", err
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
		return "", "", "", err
	}

	defer resp.Body.Close()

	// URL final setelah redirect
	finalURL := resp.Request.URL.String()

	println("FINAL URL:", finalURL)

	// baca seluruh html
	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return "", "", "", err
	}

	htmlContent := string(bodyBytes)

	// parse html
	doc, err := goquery.NewDocumentFromReader(
		strings.NewReader(htmlContent),
	)

	if err != nil {
		return "", "", "", err
	}

	// =========================
	// GOOGLE NEWS DETECT
	// =========================

	if strings.Contains(finalURL, "news.google.com") {

		// cari semua link keluar
		foundURL := ""

		doc.Find("a").Each(func(i int, s *goquery.Selection) {

			href, exists := s.Attr("href")

			if exists {

				// skip internal google
				if strings.Contains(href, "google.com") {
					return
				}

				// url asli
				if strings.HasPrefix(href, "http") {

					foundURL = href
				}
			}
		})

		// kalau ketemu url asli
		if foundURL != "" {

			println("FOUND REAL URL:", foundURL)

			return ScrapeArticle(foundURL)
		}
	}

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

	return title, content, image, nil
}
