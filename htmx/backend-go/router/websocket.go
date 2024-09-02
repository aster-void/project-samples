package router

import (
	"bytes"
	"encoding/json"
	"html/template"
	"io"
	"log"

	"github.com/aster-void/project-samples/htmx/backend-go/common"
	"github.com/aster-void/project-samples/htmx/backend-go/ws"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var tmpl = template.Must(template.New("/websocket").Parse(`
<div id="messages-container" hx-swap-oob="beforeend">
	<div class="message-box">
		{{ .Content }}
	</div>
</div>
<form id="input-form" ws-send>
	<input name="content" placeholder="message" autofocus>
    <button>送信</button>
</form>
`,
))

var hub = ws.NewHub(func(m common.Message) string {
	var b bytes.Buffer
	err := tmpl.Execute(&b, m)
	if err != nil {
		log.Println("tmpl.Execute", err.Error())
	}
	html, err := io.ReadAll(&b)
	if err != nil {
		log.Println("io.ReadAll", err.Error())
		return ""
	}
	return string(html)
})

func onRecv(db *gorm.DB) func(*ws.Client[common.Message], []byte) {
	return func(c *ws.Client[common.Message], b []byte) {
		var m common.Message
		err := json.NewDecoder(bytes.NewReader(b)).Decode(&m)
		if err != nil {
			log.Println(err)
			return
		}
		if err := db.Create(&m).Error; err != nil {
			log.Println(err)
			return
		}
		c.BroadCastToAll(m)
	}
}

func WebSocket(db *gorm.DB) func(echo.Context) error {
	return func(c echo.Context) error {
		_, err := ws.NewClient(c.Response(), c.Request(), onRecv(db), hub)
		if err != nil {
			log.Println(err)
			return c.String(500, "Failed to initialize websocket")
		}
		return nil
	}
}
