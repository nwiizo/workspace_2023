package main

import (
	"context"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

const (
	baseURL       = "http://localhost:1323"
	numberOfUsers = 100
	numberOfLoops = 10
	testDuration  = 10 * time.Second
)

var client *http.Client

func main() {
	client = &http.Client{
		Timeout: 5 * time.Second,
	}

	wg := &sync.WaitGroup{}
	ctx, cancel := context.WithTimeout(context.Background(), testDuration)
	defer cancel()

	for i := 0; i < numberOfUsers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			testUser(ctx)
		}()
	}

	wg.Wait()
}

func testUser(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			loop := rand.Intn(numberOfLoops) + 1
			for i := 0; i < loop; i++ {
				actions := []func(){
					testGetBooks,
					testGetBook,
					testAddBook,
					testUpdateBook,
					testDeleteBook,
				}
				action := actions[rand.Intn(len(actions))]
				action()
				time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
			}
		}
	}
}

func testGetBooks() {
	resp, err := client.Get(baseURL + "/books")
	if err != nil {
		fmt.Printf("Error fetching books: %v\n", err)
		return
	}
	resp.Body.Close()
}

func testGetBook() {
	// Assuming we have some book IDs available for testing
	bookID := "1"
	resp, err := client.Get(baseURL + "/books/" + bookID)
	if err != nil {
		fmt.Printf("Error fetching book: %v\n", err)
		return
	}
	resp.Body.Close()
}

// Add testAddBook, testUpdateBook, testDeleteBook functions here
func testAddBook() {
	book := `{"isbn": "1234567890", "title": "Example Book", "author": {"firstname": "John", "lastname": "Doe"}}`
	req, err := http.NewRequest(http.MethodPost, baseURL+"/books", strings.NewReader(book))
	if err != nil {
		fmt.Printf("Error creating request for adding book: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error adding book: %v\n", err)
		return
	}
	resp.Body.Close()
}

func testUpdateBook() {
	// Assuming we have some book IDs available for testing
	bookID := "1"
	updatedBook := `{"isbn": "1234567890", "title": "Updated Example Book", "author": {"firstname": "John", "lastname": "Doe"}}`
	req, err := http.NewRequest(
		http.MethodPut,
		baseURL+"/books/"+bookID,
		strings.NewReader(updatedBook),
	)
	if err != nil {
		fmt.Printf("Error creating request for updating book: %v\n", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error updating book: %v\n", err)
		return
	}
	resp.Body.Close()
}

func testDeleteBook() {
	// Assuming we have some book IDs available for testing
	bookID := "1"
	req, err := http.NewRequest(http.MethodDelete, baseURL+"/books/"+bookID, nil)
	if err != nil {
		fmt.Printf("Error creating request for deleting book: %v\n", err)
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error deleting book: %v\n", err)
		return
	}
	resp.Body.Close()
}
