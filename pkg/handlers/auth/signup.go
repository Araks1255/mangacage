package auth

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/common/utils"

	"github.com/gin-gonic/gin"
)

func (h handler) Signup(c *gin.Context) {
	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	var existingUserID uint
	h.DB.Raw("SELECT id FROM users WHERE user_name = ?", user.UserName).Scan(&existingUserID)
	if existingUserID != 0 {
		c.AbortWithStatusJSON(403, gin.H{"error": "Вы уже зарегистрированы"})
		return
	}

	var errHash error
	user.Password, errHash = utils.GenerateHashPassword(user.Password)
	if errHash != nil {
		log.Println(errHash)
		c.AbortWithStatusJSON(500, gin.H{"error": errHash.Error()})
		return
	}

	if result := h.DB.Create(&user); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(201, gin.H{"success": "Регистрация прошла успешно"})
}
