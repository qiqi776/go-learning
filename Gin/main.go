package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	// r := routers.SetupRouter()
	// if err := r.Run(":8080"); err != nil {
	// 	fmt.Println("err:", err)
	// }
	r := gin.Default()

	r.GET("/hello", func(c *gin.Context) {
		c.String(http.StatusOK, "ok")
	})

	r.Run()
}
