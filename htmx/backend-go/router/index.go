package router

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"log"
	"os"

	"github.com/aster-void/project-samples/htmx/backend-go/common"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

var index_tmpl *template.Template

func init() {
	file, err := os.ReadFile("../frontend/index.htmx.tmpl")
	if err != nil {
		log.Fatalln(err)
	}
	index_tmpl = template.Must(template.New("a").Parse(string(file)))
}

func Index(db *gorm.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		var messages []common.Message
		if err := db.Find(&messages).Error; err != nil {
			fmt.Println(err.Error())
			return c.String(500, "Failed to get messages from db")
		}
		var b bytes.Buffer
		if err := index_tmpl.Execute(&b, messages); err != nil {
			return c.String(500, "Failed to realize template")
		}
		htmx, err := io.ReadAll(&b)
		if err != nil {
			fmt.Println(err)
			return c.String(500, "Failed to read from buf")
		}
		return c.HTML(200, string(htmx))
	}
}
