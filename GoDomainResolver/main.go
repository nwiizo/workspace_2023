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

const (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	intset  = "1234567890"
	emoji   = "😀😁😂🤣😃😄😅😆😉😊😋😎😍😘😗😙😚😇🤗🤩🤔🤨😐😑😶😏😣😥😮"
	kansuji = "一二三四五六七八九十百千万億兆京垓"
)

// randNum generates a random number of the specified length
// for example, randNum(4) will generate a random number of length 4
func randNum(length int) string {
	rand.Seed(time.Now().UnixNano())
	result := make([]byte, length)
	for i := range result {
		result[i] = intset[rand.Intn(len(intset))]
	}
	return string(result)
}

// randomString generates a random string of the specified length
// for example, randomString(10) will generate a random string of length 10
func randomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

// resolveName resolves a name and prints the result
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
		randomName := randomString(10) + randNum(4) + ".com"
		wg.Add(1)
		go resolveName(randomName, &wg)
	}

	wg.Wait()
	fmt.Println("All name resolutions complete.")
}
