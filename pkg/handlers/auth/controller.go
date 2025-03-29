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

	r.POST("/signup", h.Signup)
	r.POST("/login", h.Login)
	r.POST("/logout", h.Logout)
}
