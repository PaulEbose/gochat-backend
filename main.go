package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}

func main() {
	fmt.Println("GoChat v0.0.1")
	http.ListenAndServe(":8080", nil)
}

func init() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `Go to "/ws" endpoint to access WebSocket connection.`)
	})

	http.HandleFunc("/ws", serveWs)
}

// serveWs handles our WebSocket endpoint.
func serveWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Host:", r.Host)

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		fmt.Fprintf(w, `Please use a WebSocket protocol to access connection.`)
		return
	}

	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
		}

		fmt.Println(string(data))

		if err = conn.WriteMessage(messageType, data); err != nil {
			log.Println(err)
		}
	}
}
