package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func setupDatabase() *pgx.Conn {
	if conn != nil {
		return conn
	}

	var err error
	conn, err = pgx.Connect(
		context.Background(),
		"postgres://username:password@localhost:5432/database_name",
	)
	if err != nil {
		panic(fmt.Sprintf("Unable to connect to the database: %v\n", err))
	}
	return conn
}

func setupEcho(conn *pgx.Conn) *echo.Echo {
	e := echo.New()

	// Set up routes
	e.GET("/", index)
	e.GET("/books", getBooks)
	e.GET("/books/:id", getBook)
	e.POST("/books", addBook)
	e.PUT("/books/:id", updateBook)
	e.DELETE("/books/:id", deleteBook)

	return e
}

func TestIndex(t *testing.T) {
	conn := setupDatabase()
	e := setupEcho(conn)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	err := index(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, "Hello, World!", rec.Body.String())
}

func TestAddAndGetBook(t *testing.T) {
	conn := setupDatabase()
	e := setupEcho(conn)

	author := &Author{Firstname: "John", Lastname: "Doe"}
	book := &Book{Isbn: "1234567890", Title: "Test Book", Author: author}

	// Create a new book
	reqBody, _ := json.Marshal(book)
	req := httptest.NewRequest(http.MethodPost, "/books", bytes.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	err := addBook(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var createdBook Book
	json.Unmarshal(rec.Body.Bytes(), &createdBook)

	assert.NotEmpty(t, createdBook.ID)
	assert.Equal(t, book.Isbn, createdBook.Isbn)
	assert.Equal(t, book.Title, createdBook.Title)
	assert.Equal(t, book.Author.Firstname, createdBook.Author.Firstname)
	assert.Equal(t, book.Author.Lastname, createdBook.Author.Lastname)

	// Get the created book
	req = httptest.NewRequest(http.MethodGet, fmt.Sprintf("/books/%s", createdBook.ID), nil)
	rec = httptest.NewRecorder()

	c = e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(createdBook.ID)
	err = getBook(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var fetchedBook Book
	json.Unmarshal(rec.Body.Bytes(), &fetchedBook)

	assert.Equal(t, createdBook.ID, fetchedBook.ID)
	assert.Equal(t, createdBook.Isbn, fetchedBook.Isbn)
	assert.Equal(t, createdBook.Title, fetchedBook.Title)
	assert.Equal(t, createdBook.Author.Firstname, fetchedBook.Author.Firstname)
	assert.Equal(t, createdBook.Author.Lastname, fetchedBook.Author.Lastname)
}
