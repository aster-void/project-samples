package router

import (
	"bytes"
	"html/template"
	"io"
	"log"

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

func WebSocket() func(echo.Context) error {
	return func(c echo.Context) error {
		_, err := ws.NewClient(c.Response(), c.Request(), hub)
		if err != nil {
			log.Println(err)
			return c.String(500, "Failed to initialize websocket")
		}
		return nil
	}
}
