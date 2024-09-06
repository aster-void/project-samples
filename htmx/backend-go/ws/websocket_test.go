package ws_test

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"testing"
	"time"

	"github.com/aster-void/project-samples/htmx/backend-go/ws"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/websocket"
)

func SetupServer(halt <-chan bool) {
	e := echo.New()
	hub := ws.NewHub()
	e.GET("/websocket", func(c echo.Context) error {
		_, err := ws.NewClient(c.Response(), c.Request(), hub)
		if err != nil {
			log.Fatalln(err)
		}
		return nil
	})

	go func() {
		<-halt
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := e.Shutdown(ctx)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	go func() {
		if err := e.Start(":3000"); !errors.Is(err, http.ErrServerClosed) {
			log.Fatalln(err)
		}
	}()
}

func init() {
	fmt.Println("init")
}

func TestWS(t *testing.T) {
	log.Println("Starting Server")
	assert := assert.New(t)
	halt := make(chan bool, 1)
	SetupServer(halt)
	time.Sleep(10 * time.Second)
	defer func() { halt <- true }()
	log.Println("Server startup done")

	log.Println("Sending Dial")
	ws, err := websocket.Dial("ws://localhost:3000/websocket", "", "http://localhost:3000")
	if err != nil {
		log.Fatalln(err)
	}
	defer ws.Close()

	log.Println("Sending Message")
	if err := websocket.Message.Send(ws, "Hello, World!"); err != nil {
		log.Fatalln(err)
	}

	log.Println("Receiving Message")
	var msg string
	if err := websocket.Message.Receive(ws, &msg); err != nil {
		log.Fatalln(err)
	}
	log.Println("Received Message")
	assert.Equal(msg, "Hello, World!")
}
