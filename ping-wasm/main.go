package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	const port = "8181"
	http.Handle("/", http.FileServer(http.Dir(".")))

	fmt.Printf("Starting server on http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
