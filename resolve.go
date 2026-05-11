package main

import (
	"net/http"
)

func resolveGoogleNewsURL(url string) (string, error) {

	client := &http.Client{}

	resp, err := client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	return resp.Request.URL.String(), nil
}
