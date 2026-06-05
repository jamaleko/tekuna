package main

import "strings"

func FilterPriorityLinks(items []FeedItem) []FeedItem {

 var filtered []FeedItem

 priorities := []string{
  //"nasa.gov", ada kontak informasi
  "space.com",
  "kompas.com",
  "inet.detik.com/science",
  
 // "tempo.co",
  
 // "arstechnica.com",
  //"sciencedaily.com",
 }

 for _, priority := range priorities {

  for _, item := range items {

   link := strings.ToLower(item.Link)

   if strings.Contains(link, priority) {

    filtered = append(filtered, item)
   }
  }
 }

 return filtered
}
