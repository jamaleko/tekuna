package main

import (
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

    // Ambil hasil Google News RSS
    results, err := GoogleDorkSearch(query)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    // ambil link pertama
    firstLink := ""
    resolvedResults := make([]string, len(results))

    for i, link := range results {
        resolvedLink, err := ResolveURL(link)
        if err != nil {
            resolvedLink = link // fallback ke RSS link
        }
        resolvedResults[i] = resolvedLink

        // ambil link pertama resolved
        if i == 0 {
            firstLink = resolvedLink
        }
    }

    c.JSON(http.StatusOK, gin.H{
        "first_link": firstLink,
        "results":    resolvedResults,
    })
}
