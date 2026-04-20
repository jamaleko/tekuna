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
</head>
<body>
    <h1>Tekuna News 🚀</h1>
    <p>Portal berita teknologi</p>
</body>
</html>
		`))
	})

	r.Run(":10000")
}
