package main

import (
	"fmt"
	"io"
	"net/http"
	"time"
)

func AutoPost() error {

	// ===================================
	// GOOGLE DORK
	// ===================================

	query := `((teknologi OR saintek OR sains) AND (astronomi OR antariksa OR "luar angkasa" OR satelit OR roket OR NASA OR SpaceX)) -AI -hp -smartphone after:2026-05-08`

	links, err := SearchGoogleNews(query)
	if err != nil {
		return err
	}

	if len(links) == 0 {
		return fmt.Errorf("tidak ada berita")
	}

	// ===================================
	// AMBIL LINK PERTAMA
	// ===================================

	link := links[0]

	// ===================================
	// RESOLVE GOOGLE NEWS
	// ===================================

	realURL, err := resolveGoogleNewsURL(link)
	if err != nil {
		return err
	}

	fmt.Println("REAL URL:", realURL)

	// ===================================
	// SCRAPE
	// ===================================

	article, err := ScrapeArticle(realURL)
	if err != nil {
		return err
	}

	// ===================================
	// REWRITE AI
	// ===================================

	result, err := RewriteArticle(article.Title, article.Content)
	if err != nil {
		return err
	}

	// ===================================
	// DOWNLOAD IMAGE
	// ===================================

	resp, err := http.Get(article.Image)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	imageBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// ===================================
	// UPLOAD SUPABASE
	// ===================================

	filename := createSlug(result.Title) + ".jpg"

	imageURL, err := uploadBytesToSupabase(imageBytes, filename)
	if err != nil {
		return err
	}

	// ===================================
	// SAVE DATABASE
	// ===================================

	db.Create(&Berita{
		Judul:   result.Title,
		Slug:    createSlug(result.Title),
		Isi:     result.Content,
		Gambar:  imageURL,
		Tanggal: time.Now(),
	})

	fmt.Println("AUTO POST SUCCESS")

	return nil
}
