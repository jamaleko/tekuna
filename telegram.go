package main

import (
	"fmt"
	"net/http"
	"net/url"
	"io"
    "strings"
)

func SendTelegram(link string) error {

	botToken := "8353932833:AAH5pP_d4lsrMPmDXnZ-3jQrHv-x3DKdZIM"
	chatID := "-1003887957812"

	apiURL := "https://api.telegram.org/bot" +
		botToken +
		"/sendMessage"

	values := url.Values{}
	values.Set("chat_id", chatID)
	values.Set("text", link)

	resp, err := http.PostForm(apiURL, values)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf(
			"telegram status %d",
			resp.StatusCode,
		)
	}

	return nil
}
func WaitForArticle(link string, imagePath string) {
	for {
		resp, err := http.Get(link)

		if err == nil {
			body, err := io.ReadAll(resp.Body)
			resp.Body.Close()

			if err == nil && resp.StatusCode == 200 {
				html := string(body)

				if strings.Contains(html, imagePath) {
					fmt.Println("Artikel dan gambar sudah siap")
					return
				}
			}
		}

		fmt.Println("Artikel belum siap, cek lagi 30 detik...")
		time.Sleep(30 * time.Second)
	}
}
