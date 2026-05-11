package main

import (
	"net/http"

	"github.com/PuerkitoBio/goquery"
)

func ResolveGoogleNewsURL(url string) (string, error) {

	resp, err := http.Get(url)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)

	if err != nil {
		return "", err
	}

	link, exists := doc.Find(`link[rel="canonical"]`).Attr("href")

	if exists {
		return link, nil
	}

	return url, nil
}
