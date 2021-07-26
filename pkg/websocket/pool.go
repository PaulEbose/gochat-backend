package websocket

import "fmt"

type Pool struct {
	Register   chan *Client
	Unregister chan *Client
	clients    map[*Client]bool
	Broadcast  chan Message
}

// NewPool returns a pool of different channels
// used to manage client connections.
func NewPool() *Pool {
	return &Pool{
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		Broadcast:  make(chan Message),
	}
}

// Start continually listens for messages passed in through
// any of the Pool channels and act accordingly.
func (pool *Pool) Start() {
	for {
		select {
		case client := <-pool.Register:
			pool.clients[client] = true
			fmt.Println("Size of Connection Pool:", len(pool.clients))
			for c := range pool.clients {
				c.Conn.WriteJSON(Message{Type: 1, Body: "New User Joined!", UserID: client.ID})
			}
		case client := <-pool.Unregister:
			delete(pool.clients, client)
			fmt.Println("Size of Connection Pool:", len(pool.clients))
			for c := range pool.clients {
				c.Conn.WriteJSON(Message{Type: 1, Body: "User Disconnected!", UserID: client.ID})
			}

		case msg := <-pool.Broadcast:
			for c := range pool.clients {
				c.Conn.WriteJSON(msg)
			}
		}
	}
}
