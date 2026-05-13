package main

import "strings"

func RemoveBlockedKeywords(items []FeedItem) []FeedItem {

	var filtered []FeedItem

	blocked := []string{
		"ai",
		"ps5",
		"gadget",
		"smartphone",
		"hp",
		"iphone",
		"android",
	}

	for _, item := range items {

		title :=
			strings.ToLower(item.Title)

		skip := false

		for _, word := range blocked {

			if strings.Contains(title, word) {

				skip = true
				break
			}
		}

		if !skip {
			filtered = append(filtered, item)
		}
	}

	return filtered
}
