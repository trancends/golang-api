package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Book struct {
	Id          int    `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	ReleaseYear string `json:"releaseYear"`
	Pages       int    `json:"pages"`
}

var books = []Book{
	{Id: 1, Title: "Laskar Pelangi", Author: "Andrea Hirata", ReleaseYear: "2005", Pages: 529},
	{Id: 2, Title: "Madilog", Author: "Tan Malaka", ReleaseYear: "1951", Pages: 228},
	{Id: 3, Title: "Dilan", Author: "Pidi Baiq", ReleaseYear: "2014"},
}

func main() {
	router := gin.Default()

	router.GET("/books", getAllBooks)

	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func getAllBooks(c *gin.Context) {
	c.JSON(http.StatusOK, books)
}
