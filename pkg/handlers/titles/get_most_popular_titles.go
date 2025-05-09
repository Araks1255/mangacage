package titles

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetMostPopularTitles(c *gin.Context) {
	limit := 10

	if c.Query("limit") != "" {
		var err error
		limit, err = strconv.Atoi(c.Query("limit"))
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "лимит должен быть числом"})
			return
		}
	}

	var titles []models.TitleDTO

	if err := h.DB.Raw(
		`SELECT
			t.id, t.created_at, t.name, t.description,
			a.name AS author, a.id AS author_id,
			MAX(teams.name) AS team, MAX(teams.id) AS team_id,
			ARRAY_AGG(DISTINCT g.name)::TEXT[] AS genres,
			COUNT(uvs.chapter_id) AS views
		FROM
			titles AS t
			INNER JOIN authors AS a ON a.id = t.author_id
			LEFT JOIN teams ON teams.id = t.team_id
			INNER JOIN title_genres AS tg ON tg.title_id = t.id
			INNER JOIN genres AS g ON tg.genre_id = g.id
			LEFT JOIN volumes AS v ON t.id = v.title_id
			LEFT JOIN chapters AS c ON v.id = c.volume_id
			LEFT JOIN user_viewed_chapters AS uvs ON uvs.chapter_id = c.id
		GROUP
			BY t.id, a.id
		ORDER BY
			views DESC
		LIMIT ?`, limit, // Эту махину я как-нибудь потом кэширую, но пока так
	).Scan(&titles).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(titles) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено самых читаемых тайтлов"})
		return
	}

	c.JSON(200, &titles)
}
