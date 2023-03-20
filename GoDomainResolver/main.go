package main

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"sync"
	"time"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

func resolveName(name string, wg *sync.WaitGroup) {
	defer wg.Done()

	_, err := net.LookupHost(name)
	if err != nil {
		log.Printf("Error resolving name %s: %v\n", name, err)
	} else {
		log.Printf("Name %s resolved\n", name)
	}
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: go run main.go <number_of_workers>")
	}

	numWorkers, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Error parsing number of workers: %v\n", err)
	}

	var wg sync.WaitGroup

	for i := 0; i < numWorkers; i++ {
		randomName := randomString(10) + ".com"
		wg.Add(1)
		go resolveName(randomName, &wg)
	}

	wg.Wait()
	fmt.Println("All name resolutions complete.")
}
