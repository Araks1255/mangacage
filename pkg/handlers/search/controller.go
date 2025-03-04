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

	search := r.Group("/search")
	search.GET("/:type/:query", h.Search)
}
