package services

import (
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"memo/model"
)
import "gorm.io/driver/sqlite"

var db *gorm.DB

func Start(addr string) error {
	var err error
	db, err = gorm.Open(sqlite.Open("data/data.db"))
	if err != nil {
		return err
	}

	if err := model.AutoMigrate(db); err != nil {
		return err
	}

	g := gin.Default()
	api := g.Group("/api/")
	userEndpoints(api)
	memoEndpoints(api)

	g.NoRoute(
		static.Serve("/", static.LocalFile("frontend", false)),
	)

	return g.Run(addr)
}
