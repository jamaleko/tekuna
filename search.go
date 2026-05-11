package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// Struktur response SearXNG
type SearxResponse struct {
	Results []SearxResult `json:"results"`
}

type SearxResult struct {
	Title   string `json:"title"`
	URL     string `json:"url"`
	Content string `json:"content"`
}

// Handler search
func searchHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	q := url.QueryEscape(query)
	searxURL := fmt.Sprintf("https://searx.be/search?q=%s&format=json", q)

	resp, err := http.Get(searxURL)
	if err != nil {
		http.Error(w, "failed to fetch from SearXNG: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		http.Error(w, "failed to read SearXNG response: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var data SearxResponse
	if err := json.Unmarshal(body, &data); err != nil {
		http.Error(w, "failed to parse SearXNG JSON: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}
