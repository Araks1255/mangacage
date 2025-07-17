package users

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
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

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	result := tx.Exec("UPDATE users SET visible = ? WHERE id = ? AND verificated", *requestBody.Visible, claims.ID)
	if result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(409, gin.H{"error": "ваш аккаунт ещё не прошел верификацию"})
		return
	}

	filter := bson.M{"user_id": claims.ID}
	update := bson.M{"$set": bson.M{"visible": *requestBody.Visible}}

	if _, err := h.UsersProfilePictures.UpdateOne(c.Request.Context(), filter, update); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "настройки профиля успешно изменены"})
}
