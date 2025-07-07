package moderation

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
)

func (h handler) GetMyTeamOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	moderationType := c.Query("type")

	var (
		team dto.ResponseTeamDTO
		err  error
	)

	switch moderationType {
	case "new":
		err = h.DB.Table("teams_on_moderation").Select("*").Where("creator_id = ?", claims.ID).Where("existing_id IS NULL").Scan(&team).Error

	case "edited":
		err = h.DB.Table("teams_on_moderation AS tom").
			Select("tom.*, t.name AS existing").
			Joins("INNER JOIN teams AS t ON t.id = tom.existing_id").
			Where("tom.creator_id = ?", claims.ID).
			Scan(&team).Error

	case "":
		err = h.DB.Table("teams_on_moderation AS tom").
			Select("tom.*, t.name AS existing").
			Joins("LEFT JOIN teams AS t ON tom.existing_id = t.id").
			Where("tom.creator_id = ?", claims.ID).
			Scan(&team).Error

	default:
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный тип модерации"})
		return
	}

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if team.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено вашей заявки на модерацию команды"})
		return
	}

	c.JSON(200, &team)
}
