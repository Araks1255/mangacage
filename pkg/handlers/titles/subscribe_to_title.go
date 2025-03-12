package titles

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) SubscribeToTitle(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	desiredTitle := c.Param("title")

	var desiredTitleID uint
	h.DB.Raw("SELECT id FROM titles WHERE lower(name) = lower(?)", desiredTitle).Scan(&desiredTitleID)
	if desiredTitleID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "Тайтл не найден"})
		return
	}

	if result := h.DB.Exec("INSERT INTO user_titles_subscribed_to (user_id, title_id) VALUES (?, ?)", claims.ID, desiredTitleID); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(201, gin.H{"succes": "Вы успешно подписались на тайтл"})
}
