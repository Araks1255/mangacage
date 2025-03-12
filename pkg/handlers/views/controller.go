package views

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func RegisterRoutes(db *gorm.DB, r *gin.Engine) {
	h := handler{DB: db}

	r.LoadHTMLFiles("html/reading_page.html")
	r.Static("/static", "./static")

	r.GET("/:title/:volume/:chapter", h.ShowReadingPage)
}
