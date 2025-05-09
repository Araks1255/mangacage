package moderation

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetMyProfileChangesOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var editedProfile models.UserDTO

	if err := h.DB.Raw(
		`SELECT
			uom.id, uom.created_at, uom.user_name, uom.about_yourself
		FROM
			users_on_moderation AS uom
		WHERE
			uom.existing_id = ?`,
		claims.ID,
	).Scan(&editedProfile).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if editedProfile.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено ваших изменений профиля на модерации"})
		return
	}

	c.JSON(200, &editedProfile)
}
