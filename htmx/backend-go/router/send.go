package router

import (
	"github.com/aster-void/project-samples/htmx/backend-go/common"
	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

func Send(db *gorm.DB) func(c echo.Context) error {
	return func(c echo.Context) error {
		var m common.Message
		if err := c.Bind(&m); err != nil {
			return c.String(400, "Failed to bind")
		}
		if err := db.Create(&m).Error; err != nil {
			return c.String(500, "Failed to create message")
		}

		hub.Broadcast(m.Content)

		return c.NoContent(201)
	}
}
