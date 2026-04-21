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
