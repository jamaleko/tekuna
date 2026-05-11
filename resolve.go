package main

import (
	"net/http"
)

func ResolveGoogleNewsURL(rawURL string) (string, error) {

	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return nil
		},
	}

	resp, err := client.Get(rawURL)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	finalURL := resp.Request.URL.String()

	return finalURL, nil
}
