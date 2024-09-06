package main

import (
	"errors"
	"log"
	"net/http"

	"github.com/aster-void/project-samples/htmx/backend-go/database"
	"github.com/aster-void/project-samples/htmx/backend-go/router"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	db := database.Init()

	e.GET("/", router.Index(db))
	e.GET("/websocket", router.WebSocket(db))

	if err := e.Start(":3000"); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalln(err)
	}
}
