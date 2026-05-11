package main

import (
    "net/http"

    "github.com/gin-gonic/gin"
)

// searchHandler: ambil query ?q=...
func searchHandler(c *gin.Context) {

    query := c.Query("q")

    if query == "" {
        c.JSON(http.StatusBadRequest, gin.H{
            "error": "query parameter 'q' missing",
        })
        return
    }

    // Panggil fungsi search utama
    results, err := GoogleDorkSearch(query)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{
            "error": err.Error(),
        })
        return
    }

    // ambil link pertama dan resolve ke link asli
    firstLink := ""
    if len(results) > 0 {
        resolved, err := ResolveGoogleNewsURL(results[0])
        if err == nil {
            firstLink = resolved
        } else {
            firstLink = results[0]
        }
    }

    c.JSON(http.StatusOK, gin.H{
        "first_link": firstLink,
        "results":    results,
    })
}
