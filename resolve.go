package main

import (
	"net/http"
)

// ResolveURL akan mengikuti redirect dan mengembalikan URL akhir
func ResolveURL(rawURL string) (string, error) {
    client := &http.Client{
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            // Izinkan semua redirect
            return nil
        },
    }

    resp, err := client.Get(rawURL)
    if err != nil {
        return "", err
    }
    defer resp.Body.Close()

    // URL akhir setelah redirect
    finalURL := resp.Request.URL.String()
    return finalURL, nil
}

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
