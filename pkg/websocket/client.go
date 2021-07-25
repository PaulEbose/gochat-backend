package websocket

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Pool *Pool
}

type Message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

// Read indefinitely listens for new messages
// coming through the WebSocket connection
func (client *Client) Read() {
	defer func() {
		client.Pool.Unregister <- client
		client.Conn.Close()
	}()

	for {
		messageType, data, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		message := Message{
			Type: messageType,
			Body: string(data),
		}
		client.Pool.Broadcast <- message
		fmt.Printf("Message Received: %+v\n", message)
	}
}
