package main

import (
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.Data(200, "text/html; charset=utf-8", []byte(`
<!DOCTYPE html>
<html>
<head>
    <title>Tekuna News</title>
    <script src="https://cdn.tailwindcss.com"></script>
</head>
<body class="bg-gray-100">

    <!-- Navbar -->
    <nav class="bg-blue-600 p-4 text-white text-xl font-bold">
        Tekuna
    </nav>

    <!-- Content -->
    <div class="p-6">
        <h1 class="text-3xl font-bold mb-4">Berita Teknologi</h1>

        <div class="bg-white p-4 rounded shadow">
            <h2 class="text-xl font-semibold">Website Tekuna resmi launch 🚀</h2>
            <p class="text-gray-600">Ini adalah awal dari portal berita teknologi.</p>
        </div>
    </div>

</body>
</html>
		`))
	})

	r.Run(":10000")
}
