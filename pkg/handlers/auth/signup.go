package auth

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/common/utils"

	"github.com/gin-gonic/gin"
)

func (h handler) Signup(c *gin.Context) {
	var requestBody struct {
		UserName     string `json:"userName" binding:"required"`
		Password     string `json:"password" binding:"required,min=8"`
		AboutYorself string `json:"aboutYourself"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	user := models.User{
		UserName:     requestBody.UserName,
		AboutYorself: requestBody.AboutYorself,
	}

	transaction := h.DB.Begin()

	var existingUserID uint
	transaction.Raw("SELECT id FROM users WHERE user_name = ?", user.UserName).Scan(&existingUserID)
	if existingUserID != 0 {
		transaction.Rollback()
		c.AbortWithStatusJSON(403, gin.H{"error": "Вы уже зарегистрированы"})
		return
	}

	var errHash error
	user.Password, errHash = utils.GenerateHashPassword(requestBody.Password)
	if errHash != nil {
		transaction.Rollback()
		log.Println(errHash)
		c.AbortWithStatusJSON(500, gin.H{"error": errHash.Error()})
		return
	}

	if result := transaction.Create(&user); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result := transaction.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, (SELECT id FROM roles WHERE name = 'user'))", user.ID); result.Error != nil {
		transaction.Rollback()
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	transaction.Commit()

	c.JSON(201, gin.H{"success": "Регистрация прошла успешно"})
}
