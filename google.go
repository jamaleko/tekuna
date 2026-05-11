package main

import (
	"encoding/xml"
	"net/http"
	"net/url"
)

type RSS struct {
	Channel struct {
		Items []struct {
			Title string `xml:"title"`
			Link  string `xml:"link"`
		} `xml:"item"`
	} `xml:"channel"`
}

func GoogleDorkSearch(query string) ([]string, error) {

	searchURL := "https://news.google.com/rss/search?q=" + url.QueryEscape(query) + "&hl=id&gl=ID&ceid=ID:id"

	resp, err := http.Get(searchURL)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var rss RSS

	err = xml.NewDecoder(resp.Body).Decode(&rss)
	if err != nil {
		return nil, err
	}

	var results []string

	for _, item := range rss.Channel.Items {
		results = append(results, item.Link)
	}

	return results, nil
}
