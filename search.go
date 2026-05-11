package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type SearxResponse struct {
	Results []SearxResult `json:"results"`
}

type SearxResult struct {
	Title string `json:"title"`
	URL   string `json:"url"`
}

func SearchNews(query string) ([]SearxResult, error) {

	base := "https://searx.be/search"

	params := url.Values{}
	params.Add("q", query)
	params.Add("format", "json")

	fullURL := base + "?" + params.Encode()

	client := &http.Client{}

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

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

	var result SearxResponse

	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}

	return result.Results, nil
}

func RegisterSearchRoute(r *gin.Engine) {

	r.GET("/test-search", func(c *gin.Context) {

		query := `astronomi luar angkasa NASA`

		results, err := SearchNews(query)
		if err != nil {
			c.String(500, err.Error())
			return
		}

		c.JSON(200, results)
	})
}
