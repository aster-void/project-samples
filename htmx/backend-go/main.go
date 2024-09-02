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
	e.POST("/send", router.Send(db))
	e.GET("/websocket", router.WebSocket())

	if err := e.Start(":3200"); !errors.Is(err, http.ErrServerClosed) {
		log.Fatalln(err)
	}
}
