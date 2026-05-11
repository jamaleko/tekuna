package main

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
)

type SearxResult struct {
    Title string `json:"title"`
    Url   string `json:"url"`
}

type SearxResponse struct {
    Results []SearxResult `json:"results"`
}

func main() {
    query := "astronomi OR antariksa OR NASA" // ganti sesuai keyword
    base := "https://searx.be/search"        // ganti instance SearXNG lain jika error

    // build URL
    params := url.Values{}
    params.Add("q", query)
    params.Add("format", "json") // wajib agar JSON
    searchUrl := fmt.Sprintf("%s?%s", base, params.Encode())

    // request
    req, _ := http.NewRequest("GET", searchUrl, nil)
    req.Header.Set("User-Agent", "Mozilla/5.0")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)

    var result SearxResponse
    err = json.Unmarshal(body, &result)
    if err != nil {
        log.Fatal("JSON parse error:", err)
    }

    // tampilkan hasil
    for _, r := range result.Results {
        fmt.Println("Title:", r.Title)
        fmt.Println("Link :", r.Url)
        fmt.Println("------")
    }
}
