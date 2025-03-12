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

	search := r.Group("/search/:query")

	search.GET("/titles", h.SearchTitles)
	search.GET("/volumes", h.SearchVolumes)
	search.GET("/chapters", h.SearchChapters)
	search.GET("/teams", h.SearchTeams)
	search.GET("/authors", h.SearchAuthors)
}
