package users

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
)

type profileSettings struct { // В будущем может ещё что-то появится
	Visible *bool `json:"visible" binding:"required"`
}

func (h handler) ChangeProfileSettings(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var requestBody profileSettings

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	result := h.DB.Exec("UPDATE users SET visible = ? WHERE id = ? AND verificated", *requestBody.Visible, claims.ID)

	if result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(409, gin.H{"error": "ваш аккаунт еще не прошел верификацию"})
		return
	}

	c.JSON(200, gin.H{"success": "настройки профиля успешно изменены"})
}
