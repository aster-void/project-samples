package common

import "gorm.io/gorm"

type Message struct {
	gorm.Model
	Content string `form:"content"`
}
