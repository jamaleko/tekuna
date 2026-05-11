package main

import (
	"io"
	"net/http"
	"regexp"
)

func SearchGoogle(query string) ([]string, error) {

	url := "https://www.google.com/search?q=" + query

	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// WAJIB
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

	html := string(body)

	// ambil semua url
	re := regexp.MustCompile(`https://[^"&]+`)
	matches := re.FindAllString(html, -1)

	var results []string

	for _, m := range matches {

		// skip google
		if regexp.MustCompile(`google`).MatchString(m) {
			continue
		}

		results = append(results, m)
	}

	return results, nil
}
