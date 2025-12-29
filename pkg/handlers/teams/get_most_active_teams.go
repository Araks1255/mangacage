package teams

import (
	"log"

	"github.com/gin-gonic/gin"
)

func (h handler) GetMostActiveTeams(c *gin.Context) {
	res := make([]struct {
		ID                uint   `json:"id"`
		Name              string `json:"name"`
		ChaptersPublished int64  `json:"chaptersPublished"`
	}, 10)

	err := h.DB.Raw(
		`WITH recently_published_chapters AS (
			SELECT
				COUNT(c.id) AS count, c.team_id
			FROM
				chapters AS c
			WHERE
				NOW() - c.created_at <= INTERVAL '1 month'
			GROUP BY
				c.team_id
		)
		SELECT
			t.id, t.name, rpc.count AS chapters_published
		FROM
			teams AS t
		INNER JOIN
			recently_published_chapters AS rpc ON rpc.team_id = t.id
		ORDER BY
			chapters_published DESC
		LIMIT
			20`,
	).Scan(&res).Error

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(res) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено активных команд"})
		return
	}

	c.JSON(200, &res)
}
