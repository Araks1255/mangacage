package roles

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	DB *gorm.DB
}

func RegisterRoutes(db *gorm.DB, r *gin.Engine) {
	h := handler{
		DB: db,
	}

	roles := r.Group("/api/roles")
	{
		roles.GET("/", h.GetRoles)
	}
}
