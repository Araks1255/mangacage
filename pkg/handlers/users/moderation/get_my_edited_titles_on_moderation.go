package moderation

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetMyEditedTitlesOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	limit := 10
	if c.Query("limit") != "" {
		var err error
		if limit, err = strconv.Atoi(c.Query("limit")); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный лимит"})
			return
		}
	}

	var titles []models.TitleOnModerationDTO

	if err := h.DB.Raw(
		`SELECT
			tom.id, tom.created_at, tom.name, tom.description,
			ARRAY_AGG(g.name) AS genres,
			t.name AS existing, t.id AS existing_id,
			MAX(a.name) AS author, MAX(a.id) AS author_id
		FROM
			titles_on_moderation AS tom
			INNER JOIN titles AS t ON t.id = tom.existing_id
			LEFT JOIN authors AS a ON tom.author_id = a.id
			INNER JOIN title_on_moderation_genres AS tomg ON tomg.title_on_moderation_id = tom.id
			INNER JOIN genres AS g ON g.id = tomg.genre_id
		WHERE
			tom.creator_id = ?
		GROUP BY
			tom.id, t.id
		LIMIT ?`,
		claims.ID, limit,
	).Scan(&titles).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(titles) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено ваших отредактированных тайтлов на модерации"})
		return
	}

	c.JSON(200, &titles)
}
