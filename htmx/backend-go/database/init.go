package database

import (
	"log"

	"github.com/aster-void/project-samples/htmx/backend-go/common"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Init() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	if err := db.AutoMigrate(&common.Message{}); err != nil {
		log.Fatalln(err)
	}
	return db
}
