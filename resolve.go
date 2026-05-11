package main

import (
	"io"
	"net/http"
	"regexp"
	"strings"
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

	re := regexp.MustCompile(`https://[^\s"'<>]+`)
	matches := re.FindAllString(html, -1)

	for _, m := range matches {

		// skip google
		if strings.Contains(m, "google.com") {
			continue
		}

		if strings.Contains(m, "googleusercontent") {
			continue
		}

		// skip image
		if strings.Contains(m, ".png") ||
			strings.Contains(m, ".jpg") ||
			strings.Contains(m, ".jpeg") ||
			strings.Contains(m, ".webp") {
			continue
		}

		// HARUS artikel
		if strings.Contains(m, ".com") ||
			strings.Contains(m, ".id") ||
			strings.Contains(m, ".net") {

			return m, nil
		}
	}

	return "", nil
}
