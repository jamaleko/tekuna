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
		"https://search.inetol.net/search?q=" +
			url.QueryEscape(query) +
			"&categories=news&format=json"

	client := &http.Client{}

	req, err := http.NewRequest("GET", searchURL, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")

	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req)

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
