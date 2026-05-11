package main

import (
	"encoding/xml"
	"io"
	"net/http"
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
