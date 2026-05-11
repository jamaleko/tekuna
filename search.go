package main

import (
    "errors"
    "net/http"
    "github.com/gin-gonic/gin"
)

// searchHandler: ambil query ?q=...
func searchHandler(c *gin.Context) {
    query := c.Query("q")
    if query == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "query parameter 'q' missing"})
        return
    }

    // Panggil fungsi search utama
    results, err := GoogleDorkSearch(query)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"results": results})
}

// GoogleDorkSearch: logika lama untuk search Google / SerpAPI / SearXNG
func GoogleDorkSearch(query string) ([]string, error) {
    if query == "" {
        return nil, errors.New("query kosong")
    }

    // --- BEGIN: tambahkan logika search lama kamu di sini ---
    // misal Google dork scraping, SerpAPI, atau SearXNG
    // contoh dummy:
    results := []string{
        "https://example.com/article1",
        "https://example.com/article2",
    }
    // --- END: logika search ---

    return results, nil
}
