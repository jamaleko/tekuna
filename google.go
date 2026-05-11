package main

import (
	"errors"
	"io"
	"net/http"
	"net/url"
)

func GoogleDorkSearch(query string) ([]string, error) {

	searchURL :=
		"https://priv.au/search?q=" +
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

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	// DEBUG tampilkan isi response
	return nil, errors.New(string(body))
}
