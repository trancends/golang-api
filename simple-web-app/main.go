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

	// Menampilkan semua buku
	router.GET("/books", getAllBooks)

	// Menambahkan buku
	router.POST("/books/create", createBook)

	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func getAllBooks(c *gin.Context) {
	c.JSON(http.StatusOK, books)
}

func createBook(c *gin.Context) {
	var newBook Book
	// read body request and do deserialization
	err := c.ShouldBind(&newBook)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	books = append(books, newBook)
	c.JSON(http.StatusCreated, newBook)
}
