package main

import (
	"database/sql"
	"net/http"
	"simple-web-app-with-db/config"
	"strconv"

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

	// Update buku
	router.PUT("/books", updateBookById)

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

func updateBookById(c *gin.Context) {
	id := c.Param("id")
	bookId, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}

	var updateBook Book

	if err := c.ShouldBind(&updateBook); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	}
	// update book by id
	// Retrieve the current book details from the database
	var currentBook Book
	query := "SELECT id, title, author, release_year, pages FROM mst_book WHERE id = $1"
	err = db.QueryRow(query, bookId).Scan(&currentBook.ID, &currentBook.Title, &currentBook.Author, &currentBook.ReleaseYear, &currentBook.Pages)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve book"})
		}
		return
	}

	// Update fields if they are provided in the request
	if updateBook.Title != "" {
		currentBook.Title = updateBook.Title
	}
	if updateBook.Author != "" {
		currentBook.Author = updateBook.Author
	}
	if updateBook.ReleaseYear != "" {
		currentBook.ReleaseYear = updateBook.ReleaseYear
	}
	if updateBook.Pages != 0 {
		currentBook.Pages = updateBook.Pages
	}

	// Update the book details in the database
	updateQuery := `UPDATE mst_book SET title = $2, author = $3, release_year = $4, pages = $5 WHERE id = $1`
	_, err = db.Exec(updateQuery, currentBook.ID, currentBook.Title, currentBook.Author, currentBook.ReleaseYear, currentBook.Pages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
		return
	}

	c.JSON(http.StatusOK, currentBook)
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
