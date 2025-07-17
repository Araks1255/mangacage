package moderation

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
)

func (h handler) GetMyTitleOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	titleOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id тайтла на модерации"})
		return
	}

	var result dto.ResponseTitleDTO

	err = h.DB.Raw(
		`SELECT
			tom.*, t.name AS existing, a.name AS author, aom.name AS author_on_moderation,
			ARRAY(
				SELECT DISTINCT g.name FROM genres AS g
				LEFT JOIN title_on_moderation_genres AS tomg ON tomg.genre_id = g.id
				WHERE tomg.title_on_moderation_id = tom.id			
			) AS genres,
			ARRAY(
				SELECT DISTINCT tags.name FROM tags
				LEFT JOIN title_on_moderation_tags AS tomt ON tomt.tag_id = tags.id
				WHERE tomt.title_on_moderation_id = tom.id
			) AS tags,
			ARRAY(
				SELECT DISTINCT com.volume
				FROM chapters_on_moderation AS com
				WHERE com.title_on_moderation_id = tom.id
			) AS volumes
		FROM
			titles_on_moderation AS tom
			LEFT JOIN titles AS t ON tom.existing_id = t.id
			LEFT JOIN authors AS a ON tom.author_id = a.id
			LEFT JOIN authors_on_moderation AS aom ON tom.author_on_moderation_id = aom.id
		WHERE
			tom.id = ? AND tom.creator_id = ?`,
		titleOnModerationID, claims.ID,
	).Scan(&result).Error

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if result.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден среди ваших заявок на модерацию"})
		return
	}

	c.JSON(200, &result)
}
