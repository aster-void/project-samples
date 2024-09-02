package router

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"time"

	"github.com/aster-void/project-samples/htmx/backend-go/ws"
	"github.com/labstack/echo/v4"
)

var tmpl = template.Must(template.New("/send").Parse(
	`<div>
		{{ .Content }}
	</div>`,
))

var hub = ws.NewHub(func(s string) string {
	var b bytes.Buffer
	_ = tmpl.Execute(&b, s)
	html, err := io.ReadAll(&b)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(html)
})

func init() {
	go hub.Run()
	go func() {
		hub.Broadcast("Periodic Broadcasting!")
		time.Sleep(1 * time.Second)
	}()
}

func WebSocket() func(echo.Context) error {
	return func(c echo.Context) error {
		ws, err := ws.NewClient(c.Response(), c.Request(), hub)
		if err != nil {
			return err
		}
		defer ws.Close()

		for {
			// Write
			err := ws.Send([]byte("Hello, Client!"))
			if err != nil {
				c.Logger().Error(err)
			}
			time.Sleep(1 * time.Second)
		}
	}
	//  return func(c echo.Context) error {
	// 	client, err := ws.NewClient(c.Response(), c.Request(), hub)
	// 	if err != nil {
	// 		log.Println(err)
	// 		return c.String(500, "Failed to initialize websocket")
	// 	}
	// 	client.Run()
	// 	return c.String(400, "disconnected")
	// }
}
