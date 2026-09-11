package main

import (
	"fmt"
	"net/http"
	"net/url"
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
