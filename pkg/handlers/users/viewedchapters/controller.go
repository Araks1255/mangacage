package viewedchapters

import (
	"github.com/Araks1255/mangacage/pkg/middlewares"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func RegisterRoutes(db *gorm.DB, secretKey string, r *gin.Engine) {
	h := handler{DB: db}

	viewedChapters := r.Group("/api/users/me/viewed-chapters")
	viewedChapters.Use(middlewares.Auth(secretKey))
	{
		viewedChapters.DELETE("/:id", h.DeleteViewedChapter)
		viewedChapters.DELETE("by-title/:id", h.DeleteViewedChaptersByTitle)
		viewedChapters.POST("/:id", h.CreateViewedChapter)
		viewedChapters.GET("/", h.GetViewedChapters)
	}
}
