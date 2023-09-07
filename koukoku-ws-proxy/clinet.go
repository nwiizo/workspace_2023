package main

import (
	"log"
	"time"

	"github.com/gorilla/websocket"
)

func main() {
	// サーバーに接続
	serverAddr := "ws://localhost:8181/ws"
	c, _, err := websocket.DefaultDialer.Dial(serverAddr, nil)
	if err != nil {
		log.Fatal("Failed to connect:", err)
	}
	defer c.Close()

	// サーバーへメッセージを送信
	go func() {
		for {
			time.Sleep(1 * time.Second)
			message := "Hello Server!"
			if err := c.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
				log.Println("Failed to write:", err)
				return
			}
		}
	}()

	// サーバーからメッセージを受信
	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("Failed to read:", err)
			return
		}
		log.Printf("Received: %s", message)
	}
}
