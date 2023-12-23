package main

import (
	"database/sql"
	"net/http"
	"simple-web-app-with-db/config"

	"github.com/gin-gonic/gin"
)

type Book struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Author      string `json:"author"`
	ReleaseYear string `json:"releaseYear"`
	Pages       int    `json:"pages"`
}

var db = config.ConnectDB()

func main() {
	router := gin.Default()

	// Membuat buku baru
	router.POST("/books", createBook)

	// Menampilkan semua buku
	router.GET("/books", getAllBooks)

	router.Run(":8080")
}

// Handler untuk membuat buku baru
func createBook(c *gin.Context) {
	var newBook Book
	err := c.ShouldBind(&newBook)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	query := "INSERT INTO mst_book (title, author, release_year, pages) VALUES ($1, $2, $3, $4) RETURNING id"

	var bookId int
	err = db.QueryRow(query, newBook.Title, newBook.Author, newBook.ReleaseYear, newBook.Pages).Scan(&bookId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	newBook.ID = bookId
	c.JSON(http.StatusCreated, newBook)
}

// Handler untuk menampilkan semua buku atau buku berdasarkan pencarian judul
func getAllBooks(c *gin.Context) {
	searchTitle := c.Query("title")

	query := "SELECT id,title,author,release_year,pages FROM mst_book"

	var rows *sql.Rows
	var err error

	if searchTitle != "" {
		query += " WHERE title ILIKE '%' || $1 || '%'"
		rows, err = db.Query(query, searchTitle)
	} else {
		rows, err = db.Query(query)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	defer rows.Close()

	var matchedBooks []Book
	for rows.Next() {
		var book Book
		err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.ReleaseYear, &book.Pages)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		matchedBooks = append(matchedBooks, book)
	}

	if len(matchedBooks) > 0 {
		c.JSON(http.StatusOK, matchedBooks)
	} else {
		c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
	}
}
