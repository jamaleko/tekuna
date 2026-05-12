package main

import (
 "bytes"
 "fmt"
 "io"
 "net/http"
 "os"
 "path/filepath"
 "strconv"
 "time"
)

// download + upload image
func UploadImageFromURL(
 imageURL string,
 title string,
) (string, error) {

 // ====================
 // DOWNLOAD IMAGE
 // ====================

 client := &http.Client{}

 req, err := http.NewRequest(
  "GET",
  imageURL,
  nil,
 )

 if err != nil {
  return "", err
 }

 req.Header.Set(
  "User-Agent",
  "Mozilla/5.0",
 )

 resp, err := client.Do(req)

 if err != nil {
  return "", err
 }

 defer resp.Body.Close()

 if resp.StatusCode != 200 {
  return "", fmt.Errorf("gagal download image")
 }

 // ====================
 // READ IMAGE
 // ====================

 data, err := io.ReadAll(resp.Body)

 if err != nil {
  return "", err
 }

 // ====================
 // GENERATE FILENAME
 // ====================

 ext := filepath.Ext(imageURL)

 if ext == "" {
  ext = ".jpg"
 }

 // slug dari title
 slug := createSlug(title)

 // filename SEO
 filename :=
  slug +
   "-" +
   strconv.FormatInt(
    time.Now().Unix(),
    10,
   ) +
   ext

 // ====================
 // UPLOAD SUPABASE
 // ====================

 uploadURL :=
  os.Getenv("SUPABASE_URL") +
   "/storage/v1/object/images/" +
   filename

 uploadReq, err := http.NewRequest(
  "POST",
  uploadURL,
  bytes.NewReader(data),
 )

 if err != nil {
  return "", err
 }

 uploadReq.Header.Set(
  "Authorization",
  "Bearer "+os.Getenv("SUPABASE_KEY"),
 )

 uploadReq.Header.Set(
  "Content-Type",
  "application/octet-stream",
 )

 uploadResp, err := client.Do(uploadReq)

 if err != nil {
  return "", err
 }

 defer uploadResp.Body.Close()

 if uploadResp.StatusCode != 200 &&
  uploadResp.StatusCode != 201 {

  body, _ := io.ReadAll(uploadResp.Body)

  return "", fmt.Errorf(
   "upload gagal: %s",
   string(body),
  )
 }

 // ====================
 // PUBLIC URL
 // ====================

 publicURL :=
  os.Getenv("SUPABASE_URL") +
   "/storage/v1/render/image/public/images/" +
   filename

 return publicURL, nil
}
