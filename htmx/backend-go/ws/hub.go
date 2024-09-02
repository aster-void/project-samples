package ws

import (
	"fmt"
	"log"
)

type Hub[T any] struct {
	clients   map[*Client[T]]bool
	transform func(T) string
}

// only 1 mapper is allowed
func NewHub[T any](mapper ...func(T) string) *Hub[T] {
	var transform func(T) string
	if len(mapper) < 1 {
		transform = func(t T) string { return fmt.Sprint(t) }
	} else {
		transform = mapper[0]
	}
	return &Hub[T]{
		clients:   make(map[*Client[T]]bool),
		transform: transform,
	}
}

func (h *Hub[T]) register(c *Client[T]) {
	h.clients[c] = true
}
func (h *Hub[T]) unregister(c *Client[T]) {
	delete(h.clients, c)
}

func (h *Hub[T]) Broadcast(message T) {
	fmt.Println("hub: Broadcasting", message)
	fmt.Println("to: ", len(h.clients), "clients")
	var b = []byte(h.transform(message))
	for c := range h.clients {
		err := c.Send(b)
		if err != nil {
			log.Println(err)
			c.Close()
		}
	}
}
