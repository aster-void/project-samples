package ws

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/aster-void/project-samples/htmx/backend-go/channel"
	"github.com/gorilla/websocket"
)

const (
	WRITE_WAIT       = 10 * time.Second
	PONG_WAIT        = 60 * time.Second
	PING_PERIOD      = 50 * time.Second // must be less than PONG_WAIT
	MAX_MESSAGE_SIZE = 512
	SEND_BUF_SIZE    = 50
)

var (
	newline = []byte{'\n'}
)
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Client struct {
	hubs   map[*Hub]bool
	send   chan []byte
	conn   *websocket.Conn
	closed bool
}

func NewClient(w http.ResponseWriter, r *http.Request, hubs ...*Hub) (*Client, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	var hs = make(map[*Hub]bool)
	for _, hub := range hubs {
		hs[hub] = true
	}

	conn.WriteMessage(websocket.PingMessage, nil)

	return &Client{
		hubs:   hs,
		send:   make(chan []byte, 512),
		conn:   conn,
		closed: false,
	}, nil
}

func (c *Client) Run() {
	defer c.Close()

	ctx, cancel := context.WithCancel(context.Background())
	go c.runSender(ctx, cancel)
	go c.runReceiver(ctx, cancel)

	<-ctx.Done()
}

func (c *Client) runReceiver(_ context.Context, cancel context.CancelFunc) {
	defer cancel()

	c.conn.SetReadLimit(MAX_MESSAGE_SIZE)
	if err := c.conn.SetReadDeadline(time.Now().Add(PONG_WAIT)); err != nil {
		log.Println(err)
		return
	}
	c.conn.SetPongHandler(func(string) error {
		err := c.conn.SetReadDeadline(time.Now().Add(PONG_WAIT))
		if err != nil {
			log.Println(err)
		}
		return nil
	})

	for {
		mtype, r, err := c.conn.NextReader()
		if mtype != websocket.TextMessage {
			return
		}
		if err != nil {
			return
		}
		b, err := io.ReadAll(r)
		log.Println("Read message")
		if err != nil {
			return
		}
		for hub := range c.hubs {
			hub.broadcast <- b
		}
	}
}

func (c *Client) runSender(ctx context.Context, cancel context.CancelFunc) {
	ticker := time.NewTicker(PING_PERIOD)
	defer func() {
		ticker.Stop()
		cancel()
	}()

	for {
		select {
		case <-ctx.Done():
			return
		// PING
		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(WRITE_WAIT)); err != nil {
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		// SEND
		case b, ok := <-c.send:
			if !ok {
				// channel closed
				return
			}
			err := c.conn.SetWriteDeadline(time.Now().Add(WRITE_WAIT))
			if err != nil {
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			if _, err := w.Write(b); err != nil {
				return
			}
			if err := channel.Flush(w, c.send, newline...); err != nil {
				return
			}
			if err := w.Close(); err != nil {
				return
			}
		}
	}
}

func (c *Client) Send(b []byte) error {
	select {
	case c.send <- b:
		return nil
	default:
		return errors.New("disconnected")
	}
}

func (c *Client) Join(h *Hub) {
	h.register <- c
	c.hubs[h] = true
}

func (c *Client) Leave(h *Hub) {
	h.unregister <- c
	delete(c.hubs, h)
}

func (c *Client) Close() {
	if c.closed {
		return
	}
	c.closed = true

	close(c.send)
	_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
	_ = c.conn.Close()
	for hub := range c.hubs {
		c.Leave(hub)
	}
}
