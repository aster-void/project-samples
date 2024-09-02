package ws

import "log"

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
		case b, ok := <-h.broadcast:
			if !ok {
				return
			}
			for c := range h.clients {
				err := c.Send(b)
				if err != nil {
					log.Println(err)
				}
			}
		case c := <-h.register:
			h.clients[c] = true
		case c := <-h.unregister:
			delete(h.clients, c)
		}
	}
}

func (h *Hub) Broadcast(message string) {
	h.broadcast <- []byte(h.transform(message))
}
