package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"strings"
)

type RSS struct {
	Channel struct {
		Item []struct {
			Title string `xml:"title"`
			Link  string `xml:"link"`
		} `xml:"item"`
	} `xml:"channel"`
}

func ParseRSS(url string) (*RSS, error) {

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var rss RSS

	err = xml.Unmarshal(body, &rss)

	if err != nil {
		return nil, err
	}

	return &rss, nil
}
func FilterRSS(items []struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
	}) []struct {
		Title string `xml:"title"`
		Link  string `xml:"link"`
	} {

	var filtered []struct {
		Title string `xml:"title"`
		Link  string `xml:"link"`
	}

	groupA := []string{
		"teknologi",
		"saintek",
		"sains",
	}

	groupB := []string{
		"astronomi",
		"antariksa",
		"luar angkasa",
		"satelit",
		"roket",
		"nasa",
		"spacex",
	}

	blocked := []string{
		"ai",
		"hp",
		"smartphone",
	}

	for _, item := range items {

		text := strings.ToLower(item.Title)

		matchA := false
		matchB := false
		blockedFound := false

		// cek group A
		for _, keyword := range groupA {

			if strings.Contains(text, keyword) {
				matchA = true
				break
			}
		}

		// cek group B
		for _, keyword := range groupB {

			if strings.Contains(text, keyword) {
				matchB = true
				break
			}
		}

		// cek blocked
		for _, keyword := range blocked {

			if strings.Contains(text, keyword) {
				blockedFound = true
				break
			}
		}

		// lolos filter
		if (matchA || matchB) && !blockedFound {

			filtered = append(filtered, item)
		}
	}

	return filtered
}
