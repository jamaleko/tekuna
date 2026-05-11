package main

import (
	"io"
	"net/http"
	"regexp"
)

func resolveGoogleNewsURL(url string) (string, error) {

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	html := string(body)

	// cari url asli
	re := regexp.MustCompile(`https://[^"]+`)
	match := re.FindString(html)

	return match, nil
}
