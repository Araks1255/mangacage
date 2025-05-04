package titles

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
)

func (h handler) SubscribeToTitle(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	desiredTitleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id тайтла должен быть числом"})
		return
	}

	var existingTitleID uint
	h.DB.Raw("SELECT id FROM titles WHERE id = ?", desiredTitleID).Scan(&existingTitleID)
	if existingTitleID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден"})
		return
	}

	var doesUserHaveSubscriptionToThisTitle bool
	h.DB.Raw("SELECT true FROM user_titles_subscribed_to WHERE user_id = ? AND title_id = ?", claims.ID, existingTitleID).Scan(&doesUserHaveSubscriptionToThisTitle)
	if doesUserHaveSubscriptionToThisTitle {
		c.AbortWithStatusJSON(409, gin.H{"error": "вы уже подписаны на этот тайтл"})
		return
	}

	if result := h.DB.Exec("INSERT INTO user_titles_subscribed_to (user_id, title_id) VALUES (?, ?)", claims.ID, existingTitleID); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(201, gin.H{"succes": "вы успешно подписались на тайтл"})
}
