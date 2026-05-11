package main

import (
	"encoding/json"
	"net/http"
	"net/url"
)

type SearchResult struct {
	URL string `json:"url"`
}

type SearxResponse struct {
	Results []SearchResult `json:"results"`
}

func GoogleDorkSearch(query string) ([]string, error) {

	searchURL :=
		"https://searx.be/search?q=" +
			url.QueryEscape(query) +
			"&categories=news&format=json"

	resp, err := http.Get(searchURL)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	var data SearxResponse

	err = json.NewDecoder(resp.Body).Decode(&data)

	if err != nil {
		return nil, err
	}

	var results []string

	for _, item := range data.Results {
		results = append(results, item.URL)
	}

	return results, nil
}
