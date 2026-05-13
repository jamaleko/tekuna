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
  os.Getenv("FB_PAGE_TOKEN")

 api :=
  "https://graph.facebook.com/" +
   pageID +
   "/feed"

 data := url.Values{}

 data.Set(
  "message",
  title,
 )

 data.Set(
  "link",
  link,
 )

 data.Set(
  "access_token",
  token,
 )

 resp, err :=
  http.PostForm(api, data)

 if err != nil {
  return err
 }

 defer resp.Body.Close()

 body, _ :=
  io.ReadAll(resp.Body)

 println(
  "FB RESPONSE:",
  string(body),
 )

 if resp.StatusCode != 200 {

  return fmt.Errorf(
   "facebook post gagal",
  )
 }

 return nil
}
