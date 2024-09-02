package ws

import (
	"context"
	"errors"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	// "github.com/aster-void/project-samples/htmx/backend-go/channel"
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

type Client[T any] struct {
	hubs   map[*Hub[T]]bool
	onRecv func(*Client[T], []byte)

	conn   *websocket.Conn
	m      *sync.Mutex
	ctx    context.Context
	cancel context.CancelFunc
	closed bool
}

func NewClient[T any](w http.ResponseWriter, r *http.Request, onRecv func(*Client[T], []byte), hubs ...*Hub[T]) (*Client[T], error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}

	if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
		return nil, err
	}
	log.Println("Connection established with a new client!")

	ctx, cancel := context.WithCancel(context.Background())
	c := &Client[T]{
		onRecv: onRecv,
		conn:   conn,
		m:      &sync.Mutex{},
		ctx:    ctx,
		cancel: cancel,
		closed: false,
	}

	var hs = make(map[*Hub[T]]bool)
	for _, hub := range hubs {
		hs[hub] = true
		hub.clients[c] = true
	}
	c.hubs = hs

	go c.run()
	return c, nil
}

func (c *Client[T]) run() {
	defer c.Close()

	go c.runPing()
	go c.runReceiver()
	<-c.ctx.Done()
}

func (c *Client[T]) Done() <-chan struct{} {
	return c.ctx.Done()
}

func (c *Client[T]) runReceiver() {
	defer c.cancel()

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
		if err != nil {
			return
		}
		if mtype != websocket.TextMessage {
			return
		}
		b, err := io.ReadAll(r)
		if err != nil {
			return
		}
		c.onRecv(c, b)
	}
}

func (c *Client[T]) BroadCastToAll(t T) {
	for hub := range c.hubs {
		hub.Broadcast(t)
	}
}

func (c *Client[T]) runPing() {
	ticker := time.NewTicker(PING_PERIOD)
	defer func() {
		ticker.Stop()
		c.cancel()
	}()

	for {
		select {
		case <-c.ctx.Done():
			return
		// PING
		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(WRITE_WAIT)); err != nil {
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (c *Client[T]) Send(b []byte) error {
	if c.closed {
		return errors.New("disconnected")
	}
	c.m.Lock()
	defer c.m.Unlock()
	// SEND
	err := c.conn.SetWriteDeadline(time.Now().Add(WRITE_WAIT))
	if err != nil {
		return errors.New("SetWriteDeadline fail")
	}
	w, err := c.conn.NextWriter(websocket.TextMessage)
	if err != nil {
		return errors.New("NextWriter fail")
	}
	if _, err := w.Write(b); err != nil {
		return errors.New("w.Write fail")
	}
	// if err := channel.Flush(w, c.send, newline...); err != nil {
	// 	return errors.New("channel.Flush fail")
	// }
	if err := w.Close(); err != nil {
		return errors.New("w.Close fail")
	}
	return nil
}

func (c *Client[T]) Join(h *Hub[T]) {
	h.register(c)
	c.hubs[h] = true
}

func (c *Client[T]) Leave(h *Hub[T]) {
	h.unregister(c)
	delete(c.hubs, h)
}

func (c *Client[T]) Close() {
	c.closed = true
	c.cancel()
	for hub := range c.hubs {
		c.Leave(hub)
	}

	_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
	_ = c.conn.Close()
}
