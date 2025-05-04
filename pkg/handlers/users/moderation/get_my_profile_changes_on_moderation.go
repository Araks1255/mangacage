package moderation

import (
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetMyProfileChangesOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	limit := 10
	if c.Query("limit") != "" {
		var err error
		if limit, err = strconv.Atoi(c.Query("limit")); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный лимит"})
			return
		}
	}

	var editedProfile models.UserDTO

	h.DB.Raw(
		`SELECT uom.id, uom.created_at, uom.user_name, uom.description
		FROM users_on_moderation AS uom
		WHERE uom.existing_id = ?
		LIMIT ?`,
		claims.ID, limit,
	).Scan(&editedProfile)

	c.JSON(200, &editedProfile)
}
