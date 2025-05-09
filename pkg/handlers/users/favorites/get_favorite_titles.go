package favorites

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetFavoriteTitles(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	limit := 10

	if c.Query("limit") != "" {
		var err error
		if limit, err = strconv.Atoi(c.Query("limit")); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный лимит"})
			return
		}
	}

	var titles []models.TitleDTO

	if err := h.DB.Raw(
		`SELECT
			t.id, t.created_at, t.name, t.description,
			a.name AS author, a.id AS author_id,
			MAX(teams.name) AS team, MAX(teams.id) AS team_id,
			ARRAY_AGG(g.name) AS genres
		FROM
			user_favorite_titles AS uft
			INNER JOIN titles AS t ON t.id = uft.title_id
			INNER JOIN authors AS a ON a.id = t.author_id
			LEFT JOIN teams ON t.team_id = teams.id
			INNER JOIN title_genres AS tg ON t.id = tg.title_id
			INNER JOIN genres AS g ON g.id = tg.genre_id
		WHERE
			uft.user_id = ?
		GROUP BY
			t.id, a.id
		LIMIT ?`,
		claims.ID, limit,
	).Scan(&titles).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(titles) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено ваших избранных тайтлов"})
		return
	}

	c.JSON(200, &titles)
}
