package main

import (
	"encoding/xml"
	"io"
	"net/http"
	"strings"
)

type FeedItem struct {
	Title string `xml:"title" json:"title"`
	Link  string `xml:"link" json:"link"`
	Loc   string `xml:"loc" json:"-"`
}

type RSS struct {
	Channel struct {
		Item []FeedItem `xml:"item"`
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

func FilterRSS(items []FeedItem) []FeedItem {

	var filtered []FeedItem

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
		"mars",
		"bulan",
		"galaksi",
	}

	blocked := []string{
		"ai",
		"hp",
		"smartphone",
		"iphone",
		"android",
		"ios",
	}

	for _, item := range items {

		text := strings.ToLower(
			strings.TrimSpace(item.Title),
		)

		matchA := false
		matchB := false
		blockedFound := false

		// GROUP A
		for _, keyword := range groupA {

			if strings.Contains(text, keyword) {

				matchA = true
				break
			}
		}

		// GROUP B
		for _, keyword := range groupB {

			if strings.Contains(text, keyword) {

				matchB = true
				break
			}
		}

		// BLOCKED
		for _, keyword := range blocked {

			if strings.Contains(text, keyword) {

				blockedFound = true
				break
			}
		}

		// FINAL FILTER
		if (matchA || matchB) && !blockedFound {

			filtered = append(filtered, item)
		}
	}

	return filtered
}
