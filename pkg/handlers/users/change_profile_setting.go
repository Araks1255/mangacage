package users

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
)

type profileSettings struct { // В будущем может ещё что-то появится
	Visible bool `json:"visible" binding:"required"`
}

func (h handler) ChangeProfileSettings(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var requestBody profileSettings

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := h.DB.Exec("UPDATE users SET visible = ? WHERE id = ?", requestBody.Visible, claims.ID).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": "настройки профиля успешно изменены"})
}
