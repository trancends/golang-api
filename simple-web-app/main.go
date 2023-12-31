package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

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

var users = []User{
	{Id: 1, Name: "Benedictus", Age: 23},
	{Id: 2, Name: "Jullian", Age: 23},
}

func main() {
	router := gin.Default()
	router.Use(LoggerMiddleWare)

	apiGroup := router.Group("/api")
	{

		booksGroup := apiGroup.Group("/books")
		{
			booksGroup.GET("/", getAllBooks)
			booksGroup.POST("/", createBook)
			booksGroup.GET("/:id", getBookById)
			booksGroup.PUT("/:id", updateBookById)
		}
	}

	usersGroup := router.Group("/users", getAllUsers)
	usersGroup.GET("/")

	// Middleware
	// router.Use(LoggerMiddleWare)
	//
	// // Menampilkan semua buku
	// router.GET("/books", getAllBooks)
	//
	// // Menambahkan buku
	// router.POST("/books/create", createBook)
	//
	// // Menampilhkan detail buku
	// router.GET("/books/:id", getBookById)
	//
	// // Update buku
	// router.PUT("/books/:id", updateBookById)

	err := router.Run(":8080")
	if err != nil {
		panic(err)
	}
}

func getAllUsers(c *gin.Context) {
	c.JSON(http.StatusOK, users)
}

func getAllBooks(c *gin.Context) {
	fmt.Println("Get all books")
	// without query param
	// c.JSON(http.StatusOK, books)

	title := c.Query("title")

	if title == "" {
		c.JSON(http.StatusOK, books)
		return
	}

	var matchedBooks []Book

	for _, book := range books {
		if strings.Contains(strings.ToLower(book.Title), strings.ToLower(title)) {
			matchedBooks = append(matchedBooks, book)
		}
	}

	if len(matchedBooks) > 0 {
		c.JSON(http.StatusOK, matchedBooks)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "book not found"})
	}
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

func CustomMiddleware(c *gin.Context) {
	fmt.Println("Lewat CustomMiddleware . . .")
	c.Next()
	fmt.Println("Response Lewat CustomMiddleware")
}

func LoggerMiddleWare(c *gin.Context) {
	start := time.Now()
	c.Next()
	elapsed := time.Since(start).Microseconds()
	fmt.Println("Request memakan waktu sekitar", elapsed, "micro detik")
}
