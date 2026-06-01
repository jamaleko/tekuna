package main

import (
	"net/http"
	"fmt"
	"strconv"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
	"gorm.io/gorm"
//	"gorm.io/driver/sqlite"
	"regexp"
	"strings"
	"html/template"
	"gorm.io/driver/postgres"
    "os"
	"bytes"
	"io"
	"mime/multipart"
	"time"
	"golang.org/x/crypto/bcrypt"
	"path/filepath"
	"html"
	"encoding/xml"
	"math/rand"
	"encoding/base64"
    "encoding/json"
	"net/url"
	"path"
)

type Berita struct {
	ID      uint   `gorm:"primaryKey"`
	Judul   string
	Slug    string `gorm:"unique"`
	Isi     string
	Gambar  string
	SourceLink string `gorm:"unique"`
	Tanggal time.Time // ✅ tambah ini
}

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Username string `gorm:"unique"`
	Password string // ini hash, bukan plain text
}
type URL struct {
	Loc     string `xml:"loc"`
	LastMod string `xml:"lastmod,omitempty"`
}

type URLSet struct {
    XMLName xml.Name `xml:"urlset"`
    Xmlns   string   `xml:"xmlns,attr"`
    URLs    []URL    `xml:"url"`
}
var db *gorm.DB

func main() {
	fmt.Println("STEP 1 - START")

	var err error

	// koneksi database (tanpa gcc)
	//dsn := os.Getenv("DATABASE_URL")
	/*if dsn == "" {
	    // 🔥 fallback untuk lokal
	    dsn = "host=localhost user=postgres password=123 dbname=tekuna port=5432 sslmode=disable"
	}*/
	db, err = gorm.Open(postgres.New(postgres.Config{
    DSN:                  os.Getenv("DATABASE_URL"),
    PreferSimpleProtocol: true, // 🔥 WAJIB
}), &gorm.Config{
    PrepareStmt: false,
})
	if err != nil {
	    fmt.Println("STEP 2 - DB SKIPED")
	}
	fmt.Println("STEP 2 - DB CONNECTED")

	// migrate tabel
	db.AutoMigrate(&Berita{}, &User{})

	fmt.Println("STEP 3 - DATA READY")

	// setup gin
	r := gin.Default()
	r.Static("/static", "./static")
	r.StaticFile("/favicon.ico", "./favicon.ico")
	r.StaticFile("/robots.txt", "./robots.txt")
	r.StaticFile("/sitemap2.xml", "./static/sitemap2.xml")
	r.SetFuncMap(template.FuncMap{
    	"safeHTML": func(s string) template.HTML {
        	return template.HTML(s)
	    },
	    "excerpt": func(s string) string {
		    // hapus tag HTML
		    re := regexp.MustCompile("<.*?>")
		    s = re.ReplaceAllString(s, "")
		
		    // decode HTML entity (&nbsp; &ldquo; dll)
		    s = html.UnescapeString(s)
		
		    // rapikan spasi
		    s = strings.ReplaceAll(s, "\n", " ")
		    s = strings.ReplaceAll(s, "\r", " ")
		    s = strings.ReplaceAll(s, "\u00a0", " ") // nbsp jadi spasi
		
		    // potong
		    if len(s) > 120 {
		        return s[:120] + "..."
		    }
		    return s
		},
		"insertBacaJuga": func(s string, b Berita) template.HTML {

		    bacaHTML := `<div style="background:#f5f5f5;padding:15px;margin:20px 0;border-left:4px solid #007bff;">
		        <b>Baca juga:</b><br>
		        <a href="/berita/` + b.Slug + `" style="text-decoration: none;">` + b.Judul + `</a>
		    </div>`
		
		    // cari penutup </p> kedua
		    re := regexp.MustCompile(`(?i)</p>`)
				indexes := re.FindAllStringIndex(s, -1)
				
				if len(indexes) > 0 {
				    mid := (len(indexes) + 1) / 2
				
				    // hati-hati index array (mulai dari 0)
				    if mid >= len(indexes) {
				        mid = len(indexes) - 1
				    }
				
				    pos := indexes[mid][1]
				    s = s[:pos] + bacaHTML + s[pos:]
				} else {
				    s += bacaHTML
				}
		
		    return template.HTML(s)
		},
		"relativeTime": func(t time.Time) string {
			if t.IsZero() {
				return ""
			}

			diff := time.Since(t)

			// fungsi bulan indo
			bulan := []string{
				"Januari", "Februari", "Maret", "April",
				"Mei", "Juni", "Juli", "Agustus",
				"September", "Oktober", "November", "Desember",
			}

			// kalau lebih dari 24 jam
			if diff.Hours() >= 24 {
				return fmt.Sprintf("%d %s %d",
					t.Day(),
					bulan[int(t.Month())-1],
					t.Year(),
				)
			}

			hours := int(diff.Hours())
			if hours > 0 {
				return fmt.Sprintf("%d jam lalu", hours)
			}

			minutes := int(diff.Minutes())
			if minutes > 0 {
				return fmt.Sprintf("%d menit lalu", minutes)
			}

			seconds := int(diff.Seconds())
			if seconds < 5 {
				return "baru saja"
			}

			return fmt.Sprintf("%d detik lalu", seconds)
		},
		"base": func(s string) string {
    return filepath.Base(s)
},
		"metaDesc": func(s string) string {
		    // 1. hapus tag HTML
		    re := regexp.MustCompile("<.*?>")
		    s = re.ReplaceAllString(s, "")
		
		    // 2. decode HTML entity (&nbsp; &ldquo; dll)
		    s = html.UnescapeString(s)
		
		    // 3. hapus karakter aneh (quotes dll)
		    s = strings.ReplaceAll(s, `"`, "")
		    s = strings.ReplaceAll(s, `'`, "")
		    s = strings.ReplaceAll(s, "“", "")
		    s = strings.ReplaceAll(s, "”", "")
		
		    // 4. rapikan spasi
		    s = strings.ReplaceAll(s, "\n", " ")
		    s = strings.ReplaceAll(s, "\r", " ")
		    s = strings.ReplaceAll(s, "\u00a0", " ") // nbsp
		
		    reSpace := regexp.MustCompile(`\s+`)
		    s = reSpace.ReplaceAllString(s, " ")
		
		    // 5. trim
		    s = strings.TrimSpace(s)
		
		    // 6. potong
		    if len(s) > 150 {
		        return s[:150]
		    }
		
		    return s
		},
		// 👉 TAMBAHAN UNTUK PAGINATION
	    "add": func(a, b int) int {
	        return a + b
	    },
	    "sub": func(a, b int) int {
	        return a - b
	    },

	})
	r.LoadHTMLGlob("templates/*")

	store := cookie.NewStore([]byte("secret123"))
	r.Use(sessions.Sessions("mysession", store))

	admin := r.Group("/admin")
	admin.Use(AuthRequired())
	{
	    admin.GET("", adminList)
	    admin.GET("/create", adminCreateForm)
	    admin.POST("/create", adminCreate)
	    admin.GET("/edit/:id", adminEditForm)
	    admin.POST("/edit/:id", adminEdit)
	    admin.GET("/delete/:id", adminDelete)
	}
	r.GET("/github-test", TestGithub)
	r.GET("/read-mdx", ReadMDX)
	r.GET("/login", func(c *gin.Context) {
    c.HTML(200, "login.html", nil)
	})

	r.POST("/login", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")
	
		var user User
	
		// cari user di database
		if err := db.Where("username = ?", username).First(&user).Error; err != nil {
			c.String(401, "User tidak ditemukan")
			return
		}
	
		// bandingkan password hash
		err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
		if err != nil {
			c.String(401, "Password salah")
			return
		}
	
		// login sukses
		session := sessions.Default(c)
		session.Set("user", user.Username)
		session.Save()
	
		c.Redirect(302, "/admin")
	})

	r.GET("/logout", func(c *gin.Context) {
    session := sessions.Default(c)
    session.Clear()
    session.Save()

    c.Redirect(302, "/login")
	})
	r.GET("/privacy", func(c *gin.Context) {
    c.HTML(http.StatusOK, "privacy.html", gin.H{
        "Title": "Privacy Policy - tekuna.my.id",
		"Description": "Kebijakan privasi tekuna.my.id terkait penggunaan data pengguna",
	    })
	})
	
	r.GET("/disclaimer", func(c *gin.Context) {
	    c.HTML(http.StatusOK, "disclaimer.html", gin.H{
	        "Title": "Disclaimer - tekuna.my.id",
			"Description": "Halaman disclaimer tekuna.my.id menjelaskan batasan tanggung jawab atas informasi yang disajikan di website ini.",
	    })
	})
	r.GET("/sitemap.xml", func(c *gin.Context) {
    var urls []URL

    baseURL := "https://tekuna.my.id"

    urls = append(urls, URL{Loc: baseURL + "/"})
    urls = append(urls, URL{Loc: baseURL + "/privacy"})
    urls = append(urls, URL{Loc: baseURL + "/disclaimer"})

    var berita []Berita
    db.Find(&berita)

    for _, b := range berita {
        urls = append(urls, URL{
            Loc:     baseURL + "/berita/" + b.Slug,
            LastMod: b.Tanggal.Format("2006-01-02"),
        })
    }

    c.Header("Content-Type", "application/xml")

    // 🔥 WAJIB: tulis XML header
    c.Writer.Write([]byte(xml.Header))

    // 🔥 encode manual
    encoder := xml.NewEncoder(c.Writer)
    encoder.Indent("", "  ")

    err := encoder.Encode(URLSet{
        Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
        URLs:  urls,
    })

    if err != nil {
        c.String(500, "Error generate sitemap")
        return
    }
})
	// ======================
	// ROUTES
	// ======================
    // WAJIB
r.HEAD("/", func(c *gin.Context) {
    c.Status(200)
})

r.HEAD("/berita/:slug", func(c *gin.Context) {
    c.Status(200)
})

r.HEAD("/sitemap.xml", func(c *gin.Context) {
    c.Status(200)
})

// OPSIONAL
r.HEAD("/privacy", func(c *gin.Context) {
    c.Status(200)
})

r.HEAD("/disclaimer", func(c *gin.Context) {
    c.Status(200)
})

	// homepage
	r.GET("/", func(c *gin.Context) {
		var berita []Berita
	    var total int64

	    // ambil page dari URL ?page=1
	    pageStr := c.DefaultQuery("page", "1")
	    page, _ := strconv.Atoi(pageStr)
	    if page < 1 {
	        page = 1
	    }

	    limit := 5
	    offset := (page - 1) * limit

	    // hitung total data
	    db.Model(&Berita{}).Count(&total)

	    // ambil data sesuai halaman
	    db.Order("id desc").Limit(limit).Offset(offset).Find(&berita)

	    // hitung total halaman
	    totalPage := int((total + int64(limit) - 1) / int64(limit))

	    c.HTML(http.StatusOK, "index.html", gin.H{
	        "data":      berita,
	        "page":      page,
	        "totalPage": totalPage,
			"Title": "Tekuna - Portal Berita Teknologi",
    		"Description": "Portal berita teknologi terbaru Tekuna",
			"Image": "https://sjhqjzxylogbmsshixke.supabase.co/storage/v1/object/public/images/logo.png", // ✅ WAJIB
    	})
	})

	r.GET("/berita/:slug", func(c *gin.Context) {
    slug := c.Param("slug")

    var berita Berita
    if err := db.Where("slug = ?", slug).First(&berita).Error; err != nil {

    // 🔴 kalau benar-benar tidak ditemukan → 404
    if err == gorm.ErrRecordNotFound {

        var latest []Berita
        db.Order("id desc").Limit(5).Find(&latest)

        c.HTML(404, "empty.html", gin.H{
            "Title":       "Berita tidak ditemukan",
            "Description": "Berita tidak ditemukan",
            "Message":     "Berita tidak ditemukan",
            "latest":      latest,
        })
        return
    }

    // 🟡 kalau error database / koneksi → JANGAN 404
    fmt.Println("DB ERROR:", err)

    var latest []Berita
    db.Order("id desc").Limit(5).Find(&latest)

    c.HTML(200, "empty.html", gin.H{
        "Title":       "Server sibuk",
        "Description": "Terjadi kesalahan sementara",
        "Message":     "Server sedang sibuk, coba lagi nanti",
        "latest":      latest,
    })
    return
}

    // 🔥 ambil 1 artikel random selain ini
    var bacaJuga Berita
    db.Where("id != ?", berita.ID).
        Order("RANDOM()").
        Limit(1).
        Find(&bacaJuga)
    
    var latest []Berita
	db.Where("id != ?", berita.ID).
	    Order("id desc").
	    Limit(5).
	    Find(&latest)

    c.HTML(http.StatusOK, "detail.html", gin.H{
        "data":        berita,
        "baca":        bacaJuga, // 👈 kirim ke template
        "latest":      latest, // 👈 TAMBAHAN
        "Title":       berita.Judul + " - tekuna.my.id",
        "Description": berita.Isi,
		"Image":       berita.Gambar, // ✅
	    })
	})
	r.NoRoute(func(c *gin.Context) {

	    var latest []Berita
	    db.Order("id desc").Limit(5).Find(&latest)
	
	    c.HTML(404, "empty.html", gin.H{
	        "Title": "404 - Halaman tidak ditemukan",
	        "Description": "Halaman tidak ditemukan",
	        "Message": "404 page not found",
	        "latest": latest, // 🔥 sama
	    })
	})
	fmt.Println("STEP 4 - SERVER RUNNING")
	r.GET("/test-scrape", func(c *gin.Context) {

		url := c.Query("url")
	
		if url == "" {
			c.String(400, "url kosong")
			return
		}
	
		title, content, image, err := ScrapeArticle(url)
	
		if err != nil {
			c.String(500, err.Error())
			return
		}
	
		preview := content
	
		if len(preview) > 15000 {
			preview = preview[:15000]
		}
	
		c.JSON(200, gin.H{
			"title":   title,
			"image":   image,
			"content": preview,
		})
	})
	r.GET("/token-test", func(c *gin.Context) {

	    token := os.Getenv("TOKEN_TEKUNA")
	
	    if token == "" {
	
	        c.JSON(200, gin.H{
	            "status": "TOKEN KOSONG",
	        })
	
	        return
	    }
	
	    c.JSON(200, gin.H{
	        "status": "TOKEN ADA",
	        "length": len(token),
	    })
	})
	r.GET("/api/ping", func(c *gin.Context) {
    c.JSON(200, gin.H{
        "status": "ok",
        "message": "server alive",
        "time": time.Now(),
    })
})

r.HEAD("/api/ping", func(c *gin.Context) {
    c.Status(200)
})
	r.GET("/test-ai", func(c *gin.Context) {

		url := c.Query("url")
	
		title, content, _, err := ScrapeArticle(url)
	
		if err != nil {
			c.String(500, err.Error())
			return
		}
	
		source := "judul : " + title + "\n\n" + content
	
		result, err := GenerateAIArticle(source)
	
		if err != nil {
			c.String(500, err.Error())
			return
		}
	
		c.String(200, result)
	})
	r.GET("/test-resolve", func(c *gin.Context) {

		url := c.Query("url")
	
		realURL, err := ResolveGoogleNewsURL(url)
		if err != nil {
			c.String(500, err.Error())
			return
		}
	
		c.String(200, realURL)
	})
	r.GET("/test-rss", func(c *gin.Context) {

		feeds := []string{
	
			"https://inet.detik.com/rss",
	
			"https://www.antaranews.com/rss/tekno.xml",
	
			"https://www.cnnindonesia.com/teknologi/rss",
	
			"https://www.nasa.gov/news-release/feed/",
	
			"https://www.space.com/feeds.xml",
	
			"https://feeds.arstechnica.com/arstechnica/science",
	
			"https://www.sciencedaily.com/rss/space_time.xml",
		}
	
		var allItems []FeedItem
	
		for _, feedURL := range feeds {
	
			rss, err := ParseRSS(feedURL)
	
			if err != nil {
	
				println("RSS ERROR:", feedURL)
	
				continue
			}
	
			allItems = append(
				allItems,
				rss.Channel.Item...,
			)
		}
	
		filtered := FilterRSS(allItems)
	
		c.JSON(200, gin.H{
			"total":   len(filtered),
			"results": filtered,
		})
	})
	r.GET("/test-google", func(c *gin.Context) {

		query := `((teknologi OR saintek OR sains) AND (astronomi OR antariksa OR "luar angkasa" OR satelit OR roket OR NASA OR SpaceX)) -AI -hp -smartphone`
	
		results, err := GoogleDorkSearch(query)
	
		if err != nil {
			c.String(500, err.Error())
			return
		}
	
		c.JSON(200, gin.H{
			"results": results,
		})
	})
	r.GET("/test-decode", func(c *gin.Context) {

		link := c.Query("url")
	
		decoded := DecodeGoogleNewsURL(link)
	
		c.JSON(200, gin.H{
			"decoded": decoded,
		})
	
	})
	r.GET("/test-sitemap", func(c *gin.Context) {

		sitemap, err := ParseSitemap(
			"https://www.kompas.com/sitemap-news-sains.xml",
		)
	
		if err != nil {
	
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
	
			return
		}
	
		c.JSON(200, sitemap.URLs)
	})
	r.GET("/all-feed", func(c *gin.Context) {

		var allItems []FeedItem
	
		// ====================
		// RSS FEEDS
		// ====================
	
		rssFeeds := []string{
	
			"https://inet.detik.com/rss",
			 "https://dailysocial.id/feed", 
			"https://www.antaranews.com/rss/tekno.xml",
	
			"https://www.cnnindonesia.com/teknologi/rss",
			"https://rss.tempo.co/tekno",
	
			"https://www.nasa.gov/news-release/feed/",
	
			"https://www.space.com/feeds.xml",
	
			"https://feeds.arstechnica.com/arstechnica/science",
	
			"https://www.sciencedaily.com/rss/space_time.xml",
		}
	
		for _, feed := range rssFeeds {
	
			rss, err := ParseRSS(feed)
	
			if err != nil {
	
				println("RSS ERROR:", feed)
	
				continue
			}
	
			allItems = append(
				allItems,
				rss.Channel.Item...,
			)
		}
	
		// ====================
		// SITEMAP FEEDS
		// ====================
	
		sitemapFeeds := []string{
	
			"https://www.kompas.com/sitemap-news-sains.xml",
		}
	
		for _, feed := range sitemapFeeds {
	
			sitemap, err := ParseSitemap(feed)
	
			if err != nil {
	
				println("SITEMAP ERROR:", feed)
	
				continue
			}
	
			allItems = append(
				allItems,
				sitemap.URLs...,
			)
		}
	
		// ====================
		// PRIORITY FILTER
		// ====================
		
		priority := FilterPriorityLinks(allItems)
		priority = RemoveBlockedKeywords(priority)
		c.JSON(200, gin.H{
		 "total":   len(priority),
		 "results": priority,
		})
	})
	r.GET("/process-feed", func(c *gin.Context) {

		var allItems []FeedItem
	
		// ====================
		// RSS FEEDS
		// ====================
	
		rssFeeds := []string{
	
			"https://inet.detik.com/rss",
	
			"https://www.antaranews.com/rss/tekno.xml",
	
			"https://www.cnnindonesia.com/teknologi/rss",
	
			"https://www.nasa.gov/news-release/feed/",
	
			"https://www.space.com/feeds.xml",
	
			"https://feeds.arstechnica.com/arstechnica/science",
	
			"https://www.sciencedaily.com/rss/space_time.xml",
		}
	
		for _, feed := range rssFeeds {
	
			rss, err := ParseRSS(feed)
	
			if err != nil {
				println("RSS ERROR:", feed)
				continue
			}
	
			allItems = append(
				allItems,
				rss.Channel.Item...,
			)
		}
	
		// ====================
		// SITEMAP FEEDS
		// ====================
	
		sitemapFeeds := []string{
	
			"https://www.kompas.com/sitemap-news-sains.xml",
		}
	
		for _, feed := range sitemapFeeds {
	
			sitemap, err := ParseSitemap(feed)
	
			if err != nil {
				println("SITEMAP ERROR:", feed)
				continue
			}
	
			allItems = append(
				allItems,
				sitemap.URLs...,
			)
		}
	
		// ====================
		// LIMIT TEST
		// ====================
	
		maxProcess := 1
	
		if len(allItems) < maxProcess {
			maxProcess = len(allItems)
		}
	
		var results []gin.H
		rand.Shuffle(len(allItems), func(i, j int) {
		 allItems[i], allItems[j] =
		  allItems[j], allItems[i]
		})
		for i := 0; i < maxProcess; i++ {
	
			item := allItems[i]
	
			println("PROCESS:", item.Title)
	
			// ====================
			// SCRAPE ARTICLE
			// ====================
	
			title, content, image, err := ScrapeArticle(item.Link)
	
			if err != nil {
	
				results = append(results, gin.H{
					"source_title": item.Title,
					"link":         item.Link,
					"error":        err.Error(),
				})
	
				continue
			}
	
			// skip artikel kosong
			if len(content) < 500 {
	
				results = append(results, gin.H{
					"source_title": item.Title,
					"link":         item.Link,
					"status":       "content too short",
				})
	
				continue
			}
	
			// ====================
			// AI REWRITE
			// ====================
	
			source := "Judul: " + title + "\n\n" + content
	
			rewrite, err := GenerateAIArticle(source)
	
			if err != nil {
	
				results = append(results, gin.H{
					"source_title": item.Title,
					"link":         item.Link,
					"error":        err.Error(),
				})
	
				continue
			}
	
			results = append(results, gin.H{
				"source_title": item.Title,
				"link":         item.Link,
				"title":        title,
				"image":        image,
				"rewrite":      rewrite,
			})
		}
	
		c.JSON(200, gin.H{
			"total_feed":      len(allItems),
			"processed":       maxProcess,
			"processed_result": results,
		})
	})
	r.GET("/og/:file", func(c *gin.Context) {

	file := c.Param("file")

	imageURL :=
		"https://sjhqjzxylogbmsshixke.supabase.co/storage/v1/object/public/images/" +
			file

	resp, err := http.Get(imageURL)

	if err != nil {

		c.String(500, "gagal load image")
		return
	}

	defer resp.Body.Close()

	c.Header(
		"Content-Type",
		resp.Header.Get("Content-Type"),
	)

	io.Copy(c.Writer, resp.Body)
})
	r.GET("/kom-dek", func(c *gin.Context) {

	 var allItems []FeedItem
	
	 // ====================
	 // RSS FEEDS
	 // ====================
	
	 rssFeeds := []string{
	
	  "https://inet.detik.com/rss",
	
	  "https://www.antaranews.com/rss/tekno.xml",
	
	  "https://www.cnnindonesia.com/teknologi/rss",
	
	  "https://www.nasa.gov/news-release/feed/",
	
	  "https://www.space.com/feeds.xml",
	
	  "https://feeds.arstechnica.com/arstechnica/science",
	
	  "https://www.sciencedaily.com/rss/space_time.xml",
	 }
	
	 for _, feed := range rssFeeds {
	
	  rss, err := ParseRSS(feed)
	
	  if err != nil {
	
	   println("RSS ERROR:", feed)
	
	   continue
	  }
	
	  allItems = append(
	   allItems,
	   rss.Channel.Item...,
	  )
	 }
	
	 // ====================
	 // SITEMAP FEEDS
	 // ====================
	
	 sitemapFeeds := []string{
	
	  "https://www.kompas.com/sitemap-news-sains.xml",
	 }
	
	 for _, feed := range sitemapFeeds {
	
	  sitemap, err := ParseSitemap(feed)
	
	  if err != nil {
	
	   println("SITEMAP ERROR:", feed)
	
	   continue
	  }
	
	  allItems = append(
	   allItems,
	   sitemap.URLs...,
	  )
	 }
	
	 // ====================
	 // FILTER PRIORITAS
	 // ====================
	
	 priorityItems :=
	  FilterPriorityLinks(allItems)
	
	 // kosong
	 if len(priorityItems) == 0 {
	
	  c.JSON(500, gin.H{
	   "error": "priority feed kosong",
	  })
	
	  return
	 }
	
	 // ====================
	 // RANDOM ARTICLE
	 // ====================
	
	 randIndex :=
	  rand.Intn(len(priorityItems))
	
	 item := priorityItems[randIndex]
	
	 println("RANDOM:", item.Title)
	
	 // ====================
	 // SCRAPE ARTICLE
	 // ====================
	
	 title, content, image, err :=
	  ScrapeArticle(item.Link)
	
	 if err != nil {
	
	  c.JSON(500, gin.H{
	   "error": err.Error(),
	  })
	
	  return
	 }
	
	 // content terlalu pendek
	 if len(content) < 500 {
	
	  c.JSON(500, gin.H{
	   "error": "content terlalu pendek",
	  })
	
	  return
	 }
	
	 // ====================
	 // AI REWRITE
	 // ====================
	
	 source :=
	  "Judul: " + title +
	   "\n\n" + content
	
	 rewrite, err :=
	  GenerateAIArticle(source)
	
	 if err != nil {
	
	  c.JSON(500, gin.H{
	   "error": err.Error(),
	  })
	
	  return
	 }
	
	 // ====================
	 // RESULT
	 // ====================
	
	 c.JSON(200, gin.H{
	
	  "total_feed": len(allItems),
	
	  "priority_total": len(priorityItems),
	
	  "random_source_title": item.Title,
	
	  "random_source_link": item.Link,
	
	  "title": title,
	
	  "image": image,
	
	  "rewrite": rewrite,
	 })
	})
	r.GET("/auto-post29", func(c *gin.Context) {

	 result, err := AutoPost()
	
	 if err != nil {
	
	  c.JSON(500, result)
	
	  return
	 }
	
	 c.JSON(200, result)
	})
	go StartAutoPostScheduler()
	RegisterSearchRoute(r)
	// run server
	r.StaticFile("/favicon.png", "./favicon.png")
	r.Run(":8080")
}
func StartAutoPostScheduler() {

 ticker :=
  time.NewTicker(
   30 * time.Minute,
  )

 defer ticker.Stop()

 println("AUTO POST SCHEDULER STARTED")

 for {

  select {

  case <-ticker.C:

   println("AUTO POST RUNNING")

   result, err :=
    AutoPost()

   if err != nil {

    println(
     "AUTO POST ERROR:",
     err.Error(),
    )

    continue
   }

   // aman
   title, ok :=
    result["title"].(string)

   if ok {

    println(
     "AUTO POST SUCCESS:",
     title,
    )

   } else {

    println(
     "AUTO POST FINISHED",
    )
   }
  }
 }
}
func AutoPost() (gin.H, error) {

 var allItems []FeedItem

 // ====================
 // RSS FEEDS
 // ====================

 rssFeeds := []string{

  "https://inet.detik.com/rss",
  "https://dailysocial.id/feed", 

  "https://www.antaranews.com/rss/tekno.xml",

  "https://www.cnnindonesia.com/teknologi/rss",
  "https://rss.tempo.co/tekno",

  "https://www.nasa.gov/news-release/feed/",

  "https://www.space.com/feeds.xml",

  "https://feeds.arstechnica.com/arstechnica/science",

  "https://www.sciencedaily.com/rss/space_time.xml",
 }

 for _, feed := range rssFeeds {

  rss, err := ParseRSS(feed)

  if err != nil {

   println("RSS ERROR:", feed)

   continue
  }

  allItems = append(
   allItems,
   rss.Channel.Item...,
  )
 }

 // ====================
 // SITEMAP FEEDS
 // ====================

 sitemapFeeds := []string{

  "https://www.kompas.com/sitemap-news-sains.xml",
 }

 for _, feed := range sitemapFeeds {

  sitemap, err := ParseSitemap(feed)

  if err != nil {

   println("SITEMAP ERROR:", feed)

   continue
  }

  allItems = append(
   allItems,
   sitemap.URLs...,
  )
 }

 // ====================
 // FILTER PRIORITAS
 // ====================

 priorityItems :=
  FilterPriorityLinks(allItems)
priorityItems = RemoveBlockedKeywords(priorityItems)
 if len(priorityItems) == 0 {

  return gin.H{
   "status": "priority kosong",
  }, nil
 }

// ====================
// FILTER YANG BELUM DIPOST
// TEMP: BYPASS DB
// ====================

var availableItems []FeedItem

for _, item := range priorityItems {

    if !sourceExists(item.Link) {

        availableItems = append(
            availableItems,
            item,
        )
    }
}

 // ====================
 // SEMUA SUDAH DIPOST
 // ====================

 if len(availableItems) == 0 {

  return gin.H{
   "status": "semua artikel sudah dipost",
  }, nil
 }

 // ====================
 // RANDOM ARTICLE
 // ====================

 randIndex :=
  rand.Intn(len(availableItems))

 item := availableItems[randIndex]

 println("AUTO POST:", item.Title)

 // ====================
 // SCRAPE ARTICLE
 // ====================

 title, content, image, err :=
  ScrapeArticle(item.Link)

 if err != nil {

  return gin.H{
   "error": err.Error(),
  }, err
 }

 // ====================
 // CONTENT TOO SHORT
 // ====================

 if len(content) < 500 {

  return gin.H{
   "error": "content terlalu pendek",
  }, nil
 }

 // ====================
 // AI REWRITE
 // ====================

 source :=
  "Judul: " + title +
   "\n\n" + content

 rewrite, err :=
  GenerateAIArticle(source)

 if err != nil {

  return gin.H{
   "error": err.Error(),
  }, err
 }

 // ====================
 // NORMALIZE QUOTES
 // ====================

 rewrite =
  strings.ReplaceAll(
   rewrite,
   "“",
   `"`,
  )

 rewrite =
  strings.ReplaceAll(
   rewrite,
   "”",
   `"`,
  )

 // ====================
 // PARSE TITLE + CONTENT
 // ====================

 re :=
  regexp.MustCompile(
   `(?is)JUDUL\s*:\s*(.*?)\s*ISI\s*:\s*(.*)`,
  )

 matches :=
  re.FindStringSubmatch(rewrite)

 newTitle := ""

 newContent := rewrite

 if len(matches) >= 3 {

  parsedTitle :=
   strings.TrimSpace(matches[1])

  parsedContent :=
   strings.TrimSpace(matches[2])

  // title minimal
  if len(parsedTitle) >= 10 {

   newTitle = parsedTitle
  }

  if parsedContent != "" {

   newContent = parsedContent
  }
 }

 // fallback title
 if newTitle == "" {

  newTitle = title
 }

 // ====================
 // CLEAN CONTENT
 // ====================

 newContent =
  regexp.MustCompile(
   `(?im)^JUDUL\s*:.*$`,
  ).ReplaceAllString(
   newContent,
   "",
  )

 newContent =
  regexp.MustCompile(
   `(?im)^ISI\s*:.*$`,
  ).ReplaceAllString(
   newContent,
   "",
  )

 newContent =
  strings.TrimSpace(newContent)

 // ====================
 // FORMAT HTML
 // ====================

 paragraphs :=
  strings.Split(
   newContent,
   "\n\n",
  )

 htmlContent := ""

 for _, p := range paragraphs {

  p = strings.TrimSpace(p)

  if p == "" {

   continue
  }

  // jangan bungkus heading
  if strings.HasPrefix(p, "<h2") ||
   strings.HasPrefix(p, "<h3") {

   htmlContent += p

   continue
  }
  htmlContent +=
   "<p>" + p + "</p>"
 }

 // ====================
 // UPLOAD IMAGE
 // ====================

 uploadedImage := image

 if image != "" {

  newImage, err :=
   UploadImageFromURL(
    image,
    newTitle,
   )

  if err == nil {

   uploadedImage = newImage
  }
 }

 // ====================
 // SAVE DATABASE
 // ====================

 slug := createSlug(newTitle)

 berita := Berita{

  Judul: newTitle,

  Slug: slug,

  Isi: htmlContent,

  Gambar: uploadedImage,

  SourceLink: item.Link,

  Tanggal: time.Now(),
 }
fmt.Println("SAVE TEST")

os.MkdirAll("./generated", os.ModePerm)

imageURL, _ := url.Parse(uploadedImage)

ext := path.Ext(imageURL.Path)

if ext == "" {
    ext = ".jpg"
}
loc, _ := time.LoadLocation("Asia/Jakarta")

wibTime := time.Now().In(loc)
imageName := slug + ext
imagePath := "/static/images/" + imageName
summary := createSummary(htmlContent)
mdxContent := fmt.Sprintf(`---
title: '%s'
date: "%s"
layout: PostBanner
images:
  - %s
socialImage: "%s"
summary: '%s'
source_link: "%s"
draft: false
---

%s
`,
    newTitle,
    //wibTime.Format("2006-01-02T15:04:05"),
    wibTime.Format(time.RFC3339),
	imagePath,
	imagePath,
	summary,
	item.Link,
    htmlContent,
)

filePath := "./generated/" + slug + ".mdx"

err = os.WriteFile(
    filePath,
    []byte(mdxContent),
    0644,
)

if err != nil {

    fmt.Println("GAGAL SAVE:", err)

} else {

    // DOWNLOAD IMAGE
    resp, err := http.Get(uploadedImage)
    if err != nil {

        fmt.Println("Gagal download image:", err)

    } else {

        defer resp.Body.Close()

        os.MkdirAll("./images", os.ModePerm)

        imgPath := "./images/" + imageName

        file, err := os.Create(imgPath)
        if err != nil {

            fmt.Println(err)

        } else {

            defer file.Close()

            io.Copy(file, resp.Body)

            fmt.Println("IMAGE SAVED:", imgPath)

            imageData, _ := os.ReadFile(imgPath)

            PushImageToRepo(
                imageName,
                imageData,
            )
        }
    }

    fmt.Println("MDX BERHASIL:", filePath)

    data, _ := os.ReadFile(filePath)

    PushMDXToRepo(
        slug,
        string(data),
    )
}

fmt.Println("Judul:", berita.Judul)
fmt.Println("Slug:", berita.Slug)
fmt.Println("Gambar:", berita.Gambar)
fmt.Println("Isi panjang:", len(berita.Isi))
 //db.Create(&berita)
	time.Sleep(6 * time.Minute)
	err = SendTelegram(
	"https://www.tekuna.my.id/berita/" + berita.Slug,
)

if err != nil {

	println("TELEGRAM ERROR:", err.Error())
}
	err = ShareToFacebook(
	 "https://www.tekuna.my.id/berita/" + berita.Slug,
	)
	
	if err != nil {
	
	 println("FB ERROR:", err.Error())
	}
 // ====================
 // RESULT
 // ====================

 return gin.H{

  "status": "posted",

  "title": newTitle,

  "slug": slug,

  "url": "https://tekuna.my.id/berita/" + slug,

  "image": uploadedImage,

  "source": item.Link,
 }, nil
}
// RegisterSearchRoute daftarkan route /search
func RegisterSearchRoute(r *gin.Engine) {
    r.GET("/search", searchHandler)
}
func PushImageToRepo(
    imageName string,
    imageData []byte,
) {

    content := base64.StdEncoding.EncodeToString(
        imageData,
    )

    body := map[string]interface{}{
        "message": "auto upload image",
        "content": content,
    }

    jsonData, _ := json.Marshal(body)

    req, _ := http.NewRequest(
        "PUT",
        "https://api.github.com/repos/jamaleko/tailwind-nextjs-starter-blog/contents/public/static/images/"+imageName,
        bytes.NewBuffer(jsonData),
    )

    req.Header.Set(
        "Authorization",
        "token "+os.Getenv("TOKEN_TEKUNA"),
    )

    req.Header.Set(
        "Accept",
        "application/vnd.github+json",
    )

    client := &http.Client{}

    resp, err := client.Do(req)

    if err != nil {
        fmt.Println(err)
        return
    }

    defer resp.Body.Close()

    result, _ := io.ReadAll(resp.Body)

    fmt.Println(string(result))
}
// list admin
func adminList(c *gin.Context) {
    var berita []Berita
    var total int64

    pageStr := c.DefaultQuery("page", "1")
    page, _ := strconv.Atoi(pageStr)
    if page < 1 {
        page = 1
    }

    limit := 5
    offset := (page - 1) * limit

    db.Model(&Berita{}).Count(&total)

    db.Order("id desc").Limit(limit).Offset(offset).Find(&berita)

    totalPage := int((total + int64(limit) - 1) / int64(limit))

    c.HTML(http.StatusOK, "admin_list.html", gin.H{
        "data":      berita,
        "page":      page,
        "totalPage": totalPage,
    })
}

// form tambah
func adminCreateForm(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_create.html", nil)
}

// proses tambah + upload
func PushMDXToRepo(
    slug string,
    content string,
) error {

    token := os.Getenv("TOKEN_TEKUNA")

    path := "data/blog/" + slug + ".mdx"

    body := map[string]interface{}{
        "message": "Auto post: " + slug,
        "content": base64.StdEncoding.EncodeToString(
            []byte(content),
        ),
    }

    jsonData, _ := json.Marshal(body)

    req, _ := http.NewRequest(
        "PUT",
        "https://api.github.com/repos/jamaleko/tailwind-nextjs-starter-blog/contents/"+path,
        bytes.NewBuffer(jsonData),
    )

    req.Header.Set(
        "Authorization",
        "Bearer "+token,
    )

    req.Header.Set(
        "Accept",
        "application/vnd.github+json",
    )

    client := &http.Client{}

    resp, err := client.Do(req)

    if err != nil {
        return err
    }

    defer resp.Body.Close()

    result, _ := io.ReadAll(resp.Body)

    fmt.Println(string(result))

    return nil
}
func createSummary(content string) string {

    // cari paragraf pertama
    re := regexp.MustCompile(`(?s)<p>(.*?)</p>`)
    match := re.FindStringSubmatch(content)

    var clean string

    if len(match) > 1 {

        clean = match[1]

    } else {

        // fallback kalau tidak ada <p>
        re2 := regexp.MustCompile("<[^>]*>")
        clean = re2.ReplaceAllString(content, "")
    }

    clean = strings.ReplaceAll(clean, "\n", " ")
    clean = strings.ReplaceAll(clean, "\r", " ")
    clean = strings.Join(strings.Fields(clean), " ")

	clean = strings.ReplaceAll(clean, "'", "")
	clean = strings.ReplaceAll(clean, "(", "")
	clean = strings.ReplaceAll(clean, ")", "")
	/*clean = strings.ReplaceAll(clean, `"`, "")
	clean = strings.ReplaceAll(clean, "“", "")
	clean = strings.ReplaceAll(clean, "”", "")
	clean = strings.ReplaceAll(clean, "‘", "")
	clean = strings.ReplaceAll(clean, "’", "")*/

    if len(clean) > 160 {

        cut := clean[:160]

        lastSpace := strings.LastIndex(
            cut,
            " ",
        )

        if lastSpace > 0 {
            cut = cut[:lastSpace]
        }

        clean = cut + "..."
    }

    return clean
}
func adminCreate(c *gin.Context) {
	judul := c.PostForm("judul")
	isi := c.PostForm("isi")

	file, header, _ := c.Request.FormFile("gambar")

	// bikin nama file unik
	slug := createSlug(judul)

	// ambil extension (.jpg, .png, dll)
	ext := filepath.Ext(header.Filename)
	
	// fallback kalau tidak ada extension
	if ext == "" {
	    ext = ".jpg"
	}
	
	// nama file = slug + extension
	filename := slug + "-" + strconv.FormatInt(time.Now().Unix(), 10) + ext

	// upload ke supabase
	url, err := uploadToSupabase(file, filename)
	if err != nil {
		panic(err)
	}

	// simpan ke database
	db.Create(&Berita{
		Judul:  judul,
		Slug:   createSlug(judul),
		Isi:    isi,
		Gambar: url, // ✅ pakai URL supabase
		Tanggal: time.Now(), // ✅ isi otomatis
	})

	c.Redirect(http.StatusFound, "/admin")
}

// form edit
func ReadMDX(c *gin.Context) {

 data, err := os.ReadFile(
  "./generated/dokumen-ufo-terbaru-pejabat-as-mengungkapkan-penampakan-misterius.mdx",
 )

 if err != nil {
  c.JSON(500, gin.H{
   "error": err.Error(),
  })
  return
 }

 c.JSON(200, gin.H{
  "content": string(data[:300]),
 })
}
func adminEditForm(c *gin.Context) {
	id := c.Param("id")
	var berita Berita
	db.First(&berita, id)

	c.HTML(http.StatusOK, "admin_edit.html", gin.H{
		"data": berita,
	})
}

// proses edit
func TestGithub(c *gin.Context) {

 token := os.Getenv("TOKEN_TEKUNA")

 req, _ := http.NewRequest(
  "GET",
  "https://api.github.com/user",
  nil,
 )

 req.Header.Set(
  "Authorization",
  "Bearer "+token,
 )

 client := &http.Client{}
 resp, err := client.Do(req)

 if err != nil {
  c.JSON(500, gin.H{
   "error": err.Error(),
  })
  return
 }

 defer resp.Body.Close()

 body, _ := io.ReadAll(resp.Body)

 c.Data(
  200,
  "application/json",
  body,
 )
}
func adminEdit(c *gin.Context) {
	id := c.Param("id")

	var berita Berita
	db.First(&berita, id)

	judul := c.PostForm("judul")
	isi := c.PostForm("isi")

	// generate slug baru
	slug := createSlug(judul)

	// ambil file (kalau ada)
	file, header, err := c.Request.FormFile("gambar")

	if err == nil {

		// ambil extension
		ext := strings.ToLower(filepath.Ext(header.Filename))
		if ext == "" {
			ext = ".jpg"
		}

		// 🔥 nama file pakai slug + timestamp
		filename := slug + "-" + strconv.FormatInt(time.Now().Unix(), 10) + ext

		// upload ke supabase
		url, err := uploadToSupabase(file, filename)
		if err != nil {
			panic(err)
		}

		// update gambar
		berita.Gambar = url
	}

	// update data lain
	berita.Judul = judul
	berita.Isi = isi
	berita.Slug = slug

	db.Save(&berita)

	c.Redirect(http.StatusFound, "/admin")
}

// delete
func adminDelete(c *gin.Context) {
	id := c.Param("id")
	db.Delete(&Berita{}, id)
	c.Redirect(http.StatusFound, "/admin")
}
func sourceExists(link string) bool {

    url := "https://api.github.com/repos/jamaleko/tailwind-nextjs-starter-blog/contents/data/blog"

    req, _ := http.NewRequest(
        "GET",
        url,
        nil,
    )

    req.Header.Set(
        "Authorization",
        "token "+os.Getenv("TOKEN_TEKUNA"),
    )

    client := &http.Client{}

    resp, err := client.Do(req)

    if err != nil {

        return false
    }

    defer resp.Body.Close()

    body, _ := io.ReadAll(resp.Body)

    var files []map[string]interface{}

    json.Unmarshal(
        body,
        &files,
    )

    for _, file := range files {

        downloadURL := file["download_url"].(string)

        mdxResp, err := http.Get(
            downloadURL,
        )

        if err != nil {
            continue
        }

        content, _ := io.ReadAll(
            mdxResp.Body,
        )

        mdxResp.Body.Close()

        if strings.Contains(
            string(content),
            `source_link: "`+link+`"`,
        ) {

            return true
        }
    }

    return false
}
func createSlug(text string) string {
	// lowercase
	slug := strings.ToLower(text)

	// hapus semua karakter selain huruf, angka, dan spasi
	reg := regexp.MustCompile(`[^a-z0-9\s-]`)
	slug = reg.ReplaceAllString(slug, "")

	// ganti spasi jadi dash
	slug = strings.ReplaceAll(slug, " ", "-")

	// hapus double dash
	regDash := regexp.MustCompile(`-+`)
	slug = regDash.ReplaceAllString(slug, "-")

	// trim dash di awal/akhir
	slug = strings.Trim(slug, "-")

	return slug
}

func AuthRequired() gin.HandlerFunc {
	    return func(c *gin.Context) {
	        session := sessions.Default(c)
	        user := session.Get("user")

	        if user == nil {
	            c.Redirect(http.StatusFound, "/login")
	            c.Abort()
	            return
	        }

	        c.Next()
	    }
	}
func uploadToSupabase(file multipart.File, filename string) (string, error) {
	fmt.Println("SUPABASE_URL:", os.Getenv("SUPABASE_URL"))
    fmt.Println("SUPABASE_KEY:", os.Getenv("SUPABASE_KEY"))
	url := os.Getenv("SUPABASE_URL") + "/storage/v1/object/images/" + filename

	// baca file jadi byte
	data, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+os.Getenv("SUPABASE_KEY"))
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// debug
	body, _ := io.ReadAll(resp.Body)
	fmt.Println("STATUS:", resp.Status)
	fmt.Println("RESP:", string(body))

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		return "", fmt.Errorf("upload gagal")
	}

	publicURL := os.Getenv("SUPABASE_URL") + "/storage/v1/render/image/public/images/" + filename

	return publicURL, nil
}
