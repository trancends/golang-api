package main

import (
	"net/http"
	"strconv"
	"strings"

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

	// Menampilhkan detail buku
	router.GET("/books/:id", getBookById)

	// Update buku
	router.PUT("/books/:id", updateBookById)

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

func getBookById(c *gin.Context) {
	// save value from path variable
	id := c.Param("id")

	bookId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book id"})
		return
	}

	for _, book := range books {
		if book.Id == bookId {
			c.JSON(http.StatusOK, book)
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Book not found!"})
}

func updateBookById(c *gin.Context) {
	id := c.Param("id")

	bookId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book id"})
		return
	}

	var updatedBook Book

	if err := c.ShouldBind(&updatedBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	for i, book := range books {
		if book.Id == bookId {
			if strings.TrimSpace(updatedBook.Title) != "" {
				books[i].Title = updatedBook.Title
			}
			if strings.TrimSpace(updatedBook.Author) != "" {
				books[i].Author = updatedBook.Author
			}
			if strings.TrimSpace(updatedBook.ReleaseYear) != "" {
				books[i].ReleaseYear = updatedBook.ReleaseYear
			}
			if updatedBook.Pages != 0 {
				books[i].Pages = updatedBook.Pages
			}
			c.JSON(http.StatusOK, gin.H{"message": "Book updated", "data": books[i]})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
}
