package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
)

type RSS struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Items []Item `xml:"item"`
}

type Item struct {
	Title string `xml:"title"`
	Link  string `xml:"link"`
}

func GoogleNews(query string) ([]string, error) {
	encoded := url.QueryEscape(query)

	rssURL := fmt.Sprintf(
		"https://news.google.com/rss/search?q=%s&hl=id&gl=ID&ceid=ID:id",
		encoded,
	)

	client := &http.Client{}

	req, err := http.NewRequest("GET", rssURL, nil)
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

	var rss RSS

	err = xml.Unmarshal(body, &rss)
	if err != nil {
		return nil, err
	}

	var results []string

	for _, item := range rss.Channel.Items {
		results = append(results, item.Link)
	}

	return results, nil
}

func RegisterGoogleRoute(r *gin.Engine) {
	r.GET("/test-google", func(c *gin.Context) {
		date := time.Now().AddDate(0, 0, -3).Format("2006-01-02")
		query := `((teknologi OR saintek OR sains) AND (astronomi OR antariksa OR "luar angkasa" OR satelit OR roket OR NASA OR SpaceX))` + date

		results, err := GoogleNews(query)
		if err != nil {
			c.String(500, err.Error())
			return
		}

		c.JSON(200, gin.H{
			"results": results,
		})
	})
}
