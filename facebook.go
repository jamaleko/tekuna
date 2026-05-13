package main

import (
 "bytes"
 "encoding/json"
 "fmt"
 "net/http"
 "net/url"
 "os"
)

func ShareToFacebook(
 link string,
) error {

 pageID :=
  os.Getenv("FB_PAGE_ID")

 pageToken :=
  os.Getenv("FB_PAGE_TOKEN")

 if pageID == "" ||
  pageToken == "" {

  return fmt.Errorf(
   "FB env kosong",
  )
 }

 form := url.Values{}

 // 🔥 hanya kirim link
 form.Set("link", link)

 apiURL :=
  "https://graph.facebook.com/" +
   pageID +
   "/feed?access_token=" +
   pageToken

 resp, err :=
  http.Post(
   apiURL,
   "application/x-www-form-urlencoded",
   bytes.NewBufferString(
    form.Encode(),
   ),
  )

 if err != nil {
  return err
 }

 defer resp.Body.Close()

 var result map[string]interface{}

 json.NewDecoder(resp.Body).
  Decode(&result)

 fmt.Println(
  "FB RESPONSE:",
  result,
 )

 if resp.StatusCode != 200 {

  return fmt.Errorf(
   "facebook post gagal",
  )
 }

 return nil
}
