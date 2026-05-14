package main

import (
	"net/http"
	"net/url"
)

func SendTelegram(link string) {

	botToken := "8353932833:AAH5pP_d4lsrMPmDXnZ-3jQrHv-x3DKdZIM"
	chatID := "-1003887957812"

	apiURL := "https://api.telegram.org/bot" +
		botToken +
		"/sendMessage"

	values := url.Values{}
	values.Set("chat_id", chatID)
	values.Set("text", link)

	http.PostForm(apiURL, values)
}
