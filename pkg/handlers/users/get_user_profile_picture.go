package users

import (
	"log"
	"strconv"

	"github.com/gin-gonic/gin"
)

func (h handler) GetUserProfilePicture(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id пользователя"})
		return
	}

	var path *string

	if err := h.DB.Raw("SELECT profile_picture_path FROM users WHERE id = ? AND visible", userID).Scan(&path).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if path == nil {
		c.AbortWithStatusJSON(404, gin.H{"error": "пользователь не найден"})
		return
	}

	c.File(*path)
}
