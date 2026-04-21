package main

import (
<<<<<<< HEAD
	"net/http"
	"fmt"
	"strconv"
=======
>>>>>>> 75563f31c84b570b12ecab973e536c8561cd8c91
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
    "github.com/gin-contrib/sessions/cookie"
	"gorm.io/gorm"
	"gorm.io/driver/sqlite"
	"regexp"
	"strings"
	"html/template"
)

type Berita struct {
	ID      uint   `gorm:"primaryKey"`
	Judul   string
	Slug    string `gorm:"unique"`
	Isi     string
	Gambar  string
}

var db *gorm.DB

func main() {
	fmt.Println("STEP 1 - START")

<<<<<<< HEAD
	var err error

	// koneksi database (tanpa gcc)
	db, err = gorm.Open(sqlite.Open("berita.db"), &gorm.Config{})
	if err != nil {
		panic("gagal konek database")
	}

	fmt.Println("STEP 2 - DB CONNECTED")

	// migrate tabel
	db.AutoMigrate(&Berita{})

	// seed data
	var count int64
	db.Model(&Berita{}).Count(&count)

	if count == 0 {
		db.Create(&Berita{
			Judul:  "Website Tekuna resmi launch 🚀",
			Isi:    "Portal berita teknologi mulai dibangun.",
			Gambar: "/static/images/startup.jpg",
		})

		db.Create(&Berita{
			Judul:  "AI semakin berkembang di 2026",
			Isi:    "Teknologi AI makin banyak digunakan di berbagai bidang.",
			Gambar: "/static/images/ai.jpeg",
		})
	}

	fmt.Println("STEP 3 - DATA READY")

	// setup gin
	r := gin.Default()
	r.Static("/static", "./static")
	r.SetFuncMap(template.FuncMap{
    	"safeHTML": func(s string) template.HTML {
        	return template.HTML(s)
	    },
	    "excerpt": func(s string) string {
		// hapus enter
		s = strings.ReplaceAll(s, "\n", " ")

		// potong 120 karakter
		if len(s) > 120 {
			return s[:120] + "..."
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

	    // simple hardcode dulu
	    if username == "jamaleko" && password == "badb8poefevU" {
	        session := sessions.Default(c)
	        session.Set("user", username)
	        session.Save()

	        c.Redirect(302, "/admin")
	        return
	    }

	    c.String(401, "Login gagal")
	})

	r.GET("/logout", func(c *gin.Context) {
    session := sessions.Default(c)
    session.Clear()
    session.Save()

    c.Redirect(302, "/login")
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
    	})
	})

	// detail berita
	r.GET("/berita/:slug", func(c *gin.Context) {
	slug := c.Param("slug")

	var berita Berita
	if err := db.Where("slug = ?", slug).First(&berita).Error; err != nil {
		c.String(404, "Berita tidak ditemukan")
		return
	}

	c.HTML(http.StatusOK, "detail.html", gin.H{
		"data": berita,
		"title": berita.Judul + " - tekuna.my.id",
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

	file, _ := c.FormFile("gambar")

	filename := file.Filename
	path := "static/images/" + filename

	c.SaveUploadedFile(file, path)

	db.Create(&Berita{
		Judul:  judul,
		Slug:   createSlug(judul),
		Isi:    isi,
		Gambar: "/" + path,
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

	file, _ := c.FormFile("gambar")

	if file != nil {
		filename := file.Filename
		path := "static/images/" + filename
		c.SaveUploadedFile(file, path)
		berita.Gambar = "/" + path
	}

	berita.Judul = judul
	berita.Isi = isi

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
=======
	r.GET("/", func(c *gin.Context) {
		c.Data(200, "text/html; charset=utf-8", []byte(`
<!DOCTYPE html>
<html>
<head>
    <title>Tekuna News</title>
	<meta name="description" content="Portal berita teknologi terbaru">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-100">

    <!-- Navbar -->
    <nav class="bg-blue-600 p-4 text-white text-xl font-bold">
        Tekuna
    </nav>

    <!-- Content -->
    <div class="p-4">

        <h1 class="text-2xl font-bold mb-4">Berita Teknologi</h1>

        <!-- List berita -->
        <div class="space-y-4">

            <div class="bg-white p-4 rounded shadow">
                <h2 class="text-lg font-semibold">Website Tekuna resmi launch 🚀</h2>
                <p class="text-gray-600">Portal berita teknologi mulai dibangun.</p>
            </div>

            <div class="bg-white p-4 rounded shadow">
                <h2 class="text-lg font-semibold">AI semakin berkembang di 2026 🤖</h2>
                <p class="text-gray-600">Teknologi AI makin banyak digunakan di berbagai bidang.</p>
            </div>

            <div class="bg-white p-4 rounded shadow">
                <h2 class="text-lg font-semibold">Startup Indonesia makin maju 🇮🇩</h2>
                <p class="text-gray-600">Banyak startup lokal mulai bersaing global.</p>
            </div>

        </div>

    </div>

</body>
</html>
		`))
	})

	r.Run(":10000")
}
>>>>>>> 75563f31c84b570b12ecab973e536c8561cd8c91
