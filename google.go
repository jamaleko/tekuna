package main

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func GoogleDorkSearch(query string) ([]string, error) {

	searchURL := "https://www.google.com/search?q=" + url.QueryEscape(query)

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, err
	}

	// wajib user-agent
	req.Header.Set("User-Agent", "Mozilla/5.0")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	var results []string

	doc.Find("a").Each(func(i int, s *goquery.Selection) {

		link, exists := s.Attr("href")

		if !exists {
			return
		}

		// hasil google biasanya /url?q=
		if strings.HasPrefix(link, "/url?q=") {

			link = strings.TrimPrefix(link, "/url?q=")

			parts := strings.Split(link, "&")

			if len(parts) > 0 {

				cleanURL := parts[0]

				// filter domain aneh
				if strings.HasPrefix(cleanURL, "http") {

					results = append(results, cleanURL)
				}
			}
		}
	})

	return results, nil
}
