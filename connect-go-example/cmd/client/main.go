package main

import (
	"context"
	"log"
	"net/http"

	greetv1 "example/gen/greet/v1"
	"example/gen/greet/v1/greetv1connect"

	"github.com/bufbuild/connect-go"
)

// this is a simple example of a client that uses the generated client
func main() {
	// create a new client
	// new clients can be created for each request, but it is recommended to create a single client
	client := greetv1connect.NewGreetServiceClient(
		http.DefaultClient,
		"http://localhost:8080",
	)
	res, err := client.Greet(
		context.Background(),
		connect.NewRequest(&greetv1.GreetRequest{Name: "Jane"}),
	)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println(res.Msg.Greeting)
}
