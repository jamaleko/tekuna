package main

import (
	"encoding/xml"
	"io"
	"net/http"
)

type Sitemap struct {
	URLs []FeedItem `xml:"url"`
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

	// convert <loc> -> Link
	for i := range sitemap.URLs {

		sitemap.URLs[i].Link =
			sitemap.URLs[i].Loc
	}

	return &sitemap, nil
}
