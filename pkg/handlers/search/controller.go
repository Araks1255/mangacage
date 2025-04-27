package search

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func RegisterRoutes(db *gorm.DB, r *gin.Engine) {
	h := handler{DB: db}

	r.GET("/api/search", h.Search)
}

func NewHandler(db *gorm.DB) handler {
	return handler{
		DB: db,
	}
}
