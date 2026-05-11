package main

import (
	"encoding/base64"
	"strings"
)

// DecodeGoogleNewsURL
func DecodeGoogleNewsURL(link string) string {

	// ambil bagian setelah /articles/
	parts := strings.Split(link, "/articles/")

	if len(parts) < 2 {
		return link
	}

	encoded := parts[1]

	// hapus query ?oc=5 dll
	encoded = strings.Split(encoded, "?")[0]

	// replace url-safe base64
	encoded = strings.ReplaceAll(encoded, "-", "+")
	encoded = strings.ReplaceAll(encoded, "_", "/")

	// padding
	for len(encoded)%4 != 0 {
		encoded += "="
	}

	decoded, err := base64.StdEncoding.DecodeString(encoded)

	if err != nil {
		return link
	}

	text := string(decoded)

	// cari http
	idx := strings.Index(text, "http")

	if idx == -1 {
		return link
	}

	return text[idx:]
}
