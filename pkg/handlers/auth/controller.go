package auth

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func RegisterRoutes(db *gorm.DB, client *mongo.Client, r *gin.Engine) {
	h := handler{DB: db}

	api := r.Group("/api")
	{
		api.POST("/signup", h.Signup)
		api.POST("/login", h.Login)
		api.POST("/logout", h.Logout)
	}
}
