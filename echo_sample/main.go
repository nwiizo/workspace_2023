package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// Book represents a book
type Book struct {
	ID     string  `json:"id"`
	Isbn   string  `json:"isbn"`
	Title  string  `json:"title"`
	Author *Author `json:"author"`
}

// Author represents an author
type Author struct {
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
}

// index is the handler for the route `GET /`
func index(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

var conn *pgx.Conn

// getBooks is the handler for the route `GET /books`
func getBooks(c echo.Context) error {
	rows, err := conn.Query(
		context.Background(),
		"SELECT id, isbn, title, firstname, lastname FROM books JOIN authors ON books.author_id = authors.id",
	)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Failed to get books from the database",
		)
	}
	defer rows.Close()

	books := []Book{}
	for rows.Next() {
		var b Book
		var a Author
		if err := rows.Scan(&b.ID, &b.Isbn, &b.Title, &a.Firstname, &a.Lastname); err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, "Failed to scan book data")
		}
		b.Author = &a
		books = append(books, b)
	}

	if rows.Err() != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "Failed to read book data")
	}

	return c.JSON(http.StatusOK, books)
}

// getBook is the handler for the route `GET /books/:id`
func getBook(c echo.Context) error {
	id := c.Param("id")
	var b Book
	var a Author
	err := conn.QueryRow(context.Background(), "SELECT books.id, isbn, title, firstname, lastname FROM books JOIN authors ON books.author_id = authors.id WHERE books.id=$1", id).
		Scan(&b.ID, &b.Isbn, &b.Title, &a.Firstname, &a.Lastname)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, "Book not found")
	}
	b.Author = &a
	return c.JSON(http.StatusOK, b)
}

// addBook is the handler for the route `POST /books`
func addBook(c echo.Context) error {
	book := new(Book)
	if err := c.Bind(book); err != nil {
		return err
	}
	var authorID string
	err := conn.QueryRow(context.Background(), "INSERT INTO authors (firstname, lastname) VALUES ($1, $2) RETURNING id", book.Author.Firstname, book.Author.Lastname).
		Scan(&authorID)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Failed to insert author into the database",
		)
	}
	err = conn.QueryRow(context.Background(), "INSERT INTO books (isbn, title, author_id) VALUES ($1, $2, $3) RETURNING id", book.Isbn, book.Title, authorID).
		Scan(&book.ID)
	if err != nil {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Failed to insert book into the database",
		)
	}
	return c.JSON(http.StatusCreated, book)
}

// updateBook is the handler for the route PUT /books/:id
func updateBook(c echo.Context) error {
	id := c.Param("id")
	book := new(Book)
	if err := c.Bind(book); err != nil {
		return err
	}
	commandTag, err := conn.Exec(
		context.Background(),
		"UPDATE books SET isbn=$1, title=$2 WHERE id=$3",
		book.Isbn,
		book.Title,
		id,
	)
	if err != nil || commandTag.RowsAffected() == 0 {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Failed to update book in the database",
		)
	}
	return c.JSON(http.StatusOK, book)
}

// deleteBook is the handler for the route DELETE /books/:id
func deleteBook(c echo.Context) error {
	id := c.Param("id")
	commandTag, err := conn.Exec(context.Background(), "DELETE FROM books WHERE id=$1", id)
	if err != nil || commandTag.RowsAffected() == 0 {
		return echo.NewHTTPError(
			http.StatusInternalServerError,
			"Failed to delete book from the database",
		)
	}
	return c.NoContent(http.StatusNoContent)
}

func main() {
	e := echo.New()
	// Set up middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Connect to the database
	var err error
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	conn, err = pgx.Connect(
		context.Background(),
		"postgres://username:password@localhost:5432/database_name",
	)
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to the database: %v\n", err))
	}
	defer conn.Close(context.Background())

	// Set up routes
	e.GET("/", index)
	e.GET("/books", getBooks)
	e.GET("/books/:id", getBook)
	e.POST("/books", addBook)
	e.PUT("/books/:id", updateBook)
	e.DELETE("/books/:id", deleteBook)

	// Start the server
	e.Logger.Fatal(e.Start(":1323"))
}
