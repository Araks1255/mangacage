package auth

import (
	"gorm.io/gorm"
	//"github.com/spf13/viper"
	"github.com/gin-gonic/gin"
)

type handler struct {
	DB *gorm.DB
}

func RegisterRoutes(db *gorm.DB, r *gin.Engine) {
	h := handler{DB: db}

	r.POST("/signup", h.Signup)
	r.POST("/login", h.Login)
}
