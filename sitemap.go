package main

import (
	"encoding/xml"
	"io"
	"net/http"
)

type Sitemap struct {
	URLs []URLItem `xml:"url"`
}

type URLItem struct {
	Loc string `xml:"loc"`
}

func ParseSitemap(url string) (*Sitemap, error) {

	resp, err := http.Get(url)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	var sitemap Sitemap

	err = xml.Unmarshal(body, &sitemap)

	if err != nil {
		return nil, err
	}

	return &sitemap, nil
}
