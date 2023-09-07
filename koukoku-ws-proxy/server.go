package main

import (
	"crypto/tls"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

const (
	listenPort    = ":8181"
	proxyEndpoint = "/ws"
	serverHost    = "koukoku.shadan.open.ad.jp"
	serverPort    = "992"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func main() {
	http.HandleFunc(proxyEndpoint, wsHandler)
	log.Printf("Listening on port %s", listenPort)
	log.Fatal(http.ListenAndServe(listenPort, nil))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade Error:", err)
		return
	}
	defer conn.Close()

	config := &tls.Config{
		InsecureSkipVerify: true,
	}

	tlsConn, err := tls.Dial("tcp", serverHost+":"+serverPort, config)
	if err != nil {
		log.Println("Failed to connect with koukoku:", err)
		return
	}
	defer tlsConn.Close()

	log.Println("Connected to koukoku")

	go connectRead(tlsConn, conn)
	connectWrite(conn, tlsConn)
}

func connectRead(extConn *tls.Conn, wsConn *websocket.Conn) {
	for {
		message := make([]byte, 1024)
		n, err := extConn.Read(message)
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		msg := string(message[:n])
		if strings.HasPrefix(msg, ">>") {
			err = wsConn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				log.Println("Write error:", err)
				break
			}
		}
	}
}

func connectWrite(wsConn *websocket.Conn, extConn *tls.Conn) {
	for {
		messageType, p, err := wsConn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		if messageType == websocket.TextMessage {
			_, err = extConn.Write(p)
			if err != nil {
				log.Println("Write error:", err)
				break
			}
		}
	}
}
