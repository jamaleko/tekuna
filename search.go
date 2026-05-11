package main

import (
    //"errors"
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

