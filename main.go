package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/paulebose/gochat/pkg/websocket"
)

// serveWs handles our WebSocket endpoint.
func serveWs(pool *websocket.Pool, w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Upgrade(w, r)
	if err != nil {
		fmt.Fprintf(w, `Please use a WebSocket protocol to access connection.`)
		return
	}

	client := websocket.Client{
		ID:   strings.Replace(uuid.New().String(), "-", "", -1),
		Conn: conn,
		Pool: pool,
	}
	pool.Register <- &client
	client.Read()
}

func init() {
	p := websocket.NewPool()
	go p.Start()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, `Go to "/ws" endpoint to access WebSocket connection.`)
	})

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(p, w, r)
	})
}

func main() {
	fmt.Println("GoChat v0.0.1")
	http.ListenAndServe(":8080", nil)
}
