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
)

type Berita struct {
	ID      uint   `gorm:"primaryKey"`
	Judul   string
	Slug    string `gorm:"unique"`
	Isi     string
	Gambar  string
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
	XMLName xmlns `xml:"urlset"`
	Xmlns   string `xml:"xmlns,attr"`
	URLs    []URL `xml:"url"`
}

type xmlns struct{}
var db *gorm.DB

func main() {
	fmt.Println("STEP 1 - START")

	var err error

	// koneksi database (tanpa gcc)
	dsn := os.Getenv("DATABASE_URL")

	db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
	    panic("gagal konek database")
	}
	fmt.Println("STEP 2 - DB CONNECTED")

	// migrate tabel
	db.AutoMigrate(&Berita{}, &User{})

	fmt.Println("STEP 3 - DATA READY")

	// setup gin
	r := gin.Default()
	r.Static("/static", "./static")
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
		"metaDesc": func(s string) string {
		    re := regexp.MustCompile("<.*?>")
		    s = re.ReplaceAllString(s, "")
		
		    s = strings.ReplaceAll(s, "\n", " ")
		    s = strings.ReplaceAll(s, "\r", " ")
		
		    if len(s) > 120 {
		        return s[:120]
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
	
		// 🔹 halaman statis
		urls = append(urls, URL{Loc: baseURL + "/"})
		urls = append(urls, URL{Loc: baseURL + "/privacy"})
		urls = append(urls, URL{Loc: baseURL + "/disclaimer"})
	
		// 🔹 ambil semua berita dari DB
		var berita []Berita
		db.Find(&berita)
	
		for _, b := range berita {
			urls = append(urls, URL{
				Loc:     baseURL + "/berita/" + b.Slug,
				LastMod: time.Now().Format("2006-01-02"),
			})
		}
	
		// 🔹 buat XML
		c.Header("Content-Type", "application/xml")
		c.XML(200, URLSet{
			Xmlns: "http://www.sitemaps.org/schemas/sitemap/0.9",
			URLs:  urls,
		})
	})
	// ======================
	// ROUTES
	// ======================

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
    	})
	})

	r.GET("/berita/:slug", func(c *gin.Context) {
    slug := c.Param("slug")

    var berita Berita
    if err := db.Where("slug = ?", slug).First(&berita).Error; err != nil {
        c.String(404, "Berita tidak ditemukan")
        return
    }

    // 🔥 ambil 1 artikel random selain ini
    var bacaJuga Berita
    db.Where("id != ?", berita.ID).
        Order("RANDOM()").
        Limit(1).
        Find(&bacaJuga)

    c.HTML(http.StatusOK, "detail.html", gin.H{
        "data":        berita,
        "baca":        bacaJuga, // 👈 kirim ke template
        "Title":       berita.Judul + " - tekuna.my.id",
        "Description": berita.Isi,
	    })
	})

	fmt.Println("STEP 4 - SERVER RUNNING")

	// run server
	r.StaticFile("/favicon.png", "./favicon.png")
	r.Run(":8080")
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
	})

	c.Redirect(http.StatusFound, "/admin")
}

// form edit
func adminEditForm(c *gin.Context) {
	id := c.Param("id")
	var berita Berita
	db.First(&berita, id)

	c.HTML(http.StatusOK, "admin_edit.html", gin.H{
		"data": berita,
	})
}

// proses edit
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
