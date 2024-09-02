package ws

import (
	"fmt"
	"log"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	transform  func(string) string
}

// only 1 mapper is allowed
func NewHub(mapper ...func(string) string) *Hub {
	var transform func(string) string
	if len(mapper) < 1 {
		transform = func(s string) string { return s }
	} else {
		transform = mapper[0]
	}
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte, 1),
		register:   make(chan *Client, 1),
		unregister: make(chan *Client, 1),
		transform:  transform,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case c := <-h.register:
			h.clients[c] = true
		case c := <-h.unregister:
			delete(h.clients, c)
		}
	}
}

func (h *Hub) Broadcast(message string) {
	fmt.Println("hub: Broadcasting", message)
	fmt.Println("to: ", len(h.clients), "clients")
	var b = []byte(message)
	for c := range h.clients {
		err := c.Send(b)
		if err != nil {
			log.Println(h.clients, err)
			c.Close()
		}
	}
}
