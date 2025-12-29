package titles

import (
	"fmt"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
)

func (h handler) GetSimilarTitles(c *gin.Context) {
	baseTitleId, err := strconv.ParseUint(c.Query("to"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id тайтла"})
		return
	}

	limit := uint64(10)
	if c.Query("limit") != "" {
		if limit, err = strconv.ParseUint(c.Query("limit"), 10, 64); err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный лимит"})
			return
		}
	}

	var res []dto.ResponseTitleDTO

	query := fmt.Sprintf(
		`WITH base_title_genres AS (
			SELECT genre_id FROM title_genres WHERE title_id = ?
		),
		base_title_tags AS (
			SELECT tag_id FROM title_tags WHERE title_id = ?
		),
		base_title_meta AS (
			SELECT author_id, type FROM titles WHERE id = ?
		),
		genre_scores AS (
			SELECT
				tg.title_id,
				COUNT(*) as score
			FROM
				title_genres tg
			INNER JOIN
				base_title_genres btg ON tg.genre_id = btg.genre_id
			GROUP BY
				tg.title_id
		),
		tag_scores AS (
			SELECT
				tt.title_id,
				COUNT(*) as score
			FROM
				title_tags tt
			INNER JOIN
				base_title_tags btt ON tt.tag_id = btt.tag_id
			GROUP BY
				tt.title_id
		)
		SELECT
			t.id,
			t.name
		FROM
			titles AS t
		CROSS JOIN
			base_title_meta
		LEFT JOIN
			genre_scores AS gs ON t.id = gs.title_id
		LEFT JOIN
			tag_scores AS ts ON t.id = ts.title_id
		WHERE
			t.id != ? AND NOT t.hidden
		ORDER BY
			COALESCE(gs.score, 0) +
			COALESCE(ts.score, 0) +
			(CASE WHEN t.author_id = base_title_meta.author_id THEN 10 ELSE 0 END) +
			(CASE WHEN t.type = base_title_meta.type THEN 2 ELSE 0 END) DESC
		LIMIT
			%d`,
		limit,
	)

	if err := h.DB.Raw(query, baseTitleId, baseTitleId, baseTitleId, baseTitleId).Scan(&res).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, &res)
}
