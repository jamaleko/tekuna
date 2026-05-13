package main

import (
 "bytes"
 "encoding/json"
 "fmt"
 "net/http"
 "os"
)

func ShareToFacebook(
 title string,
 link string,
) error {

 pageID :=
  os.Getenv("FB_PAGE_ID")

 token :=
  os.Getenv("FB_ACCESS_TOKEN")

 url :=
  "https://graph.facebook.com/" +
   pageID +
   "/feed"

 payload := map[string]string{
  "message": title,
  "link":    link,
  "access_token": token,
 }

 jsonData, _ :=
  json.Marshal(payload)

 resp, err :=
  http.Post(
   url,
   "application/json",
   bytes.NewBuffer(jsonData),
  )

 if err != nil {
  return err
 }

 defer resp.Body.Close()

 if resp.StatusCode != 200 {

  return fmt.Errorf(
   "facebook post gagal",
  )
 }

 return nil
}
