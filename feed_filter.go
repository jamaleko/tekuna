package main

import "strings"

func FilterPriorityLinks(items []FeedItem) []FeedItem {

 var filtered []FeedItem

 priorities := []string{
  "antaranews.com",
  "cnnindonesia.com",
  "nasa.gov",
  "kompas.com",
  "inet.detik.com/science",
  
 // "tempo.co",
  
 // "space.com",
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
