package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Halo, Enigmanians")
	})

	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}
