package viewedchapters

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
)

type getViewedChaptersParams struct {
	Limit int `form:"limit,default=20"`
	Page  int `form:"page,default=1"`
}

func (h handler) GetViewedChapters(c *gin.Context) { // Это получение конкретно истории просмотров. Выбирается самая "недавно прочитанная" глава, находится её тайтл, удаляются все остальные главы из этого тайтла. Так для второй (когда дубликаты уже удалены), третьей и т.д. (не то же самое, что GetChapters с флагом viewed=true)
	claims := c.MustGet("claims").(*auth.Claims)

	var params getViewedChaptersParams

	if err := c.ShouldBindQuery(&params); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	offset := (params.Page - 1) * params.Limit
	if offset < 0 {
		offset = 0
	}

	var chapters []dto.ResponseChapterDTO

	err := h.DB.Raw(
		`SELECT * FROM (
			SELECT
				DISTINCT ON (t.id)
				c.*, t.id AS title_id, t.name AS title, teams.name AS team, uvc.created_at AS was_read_at
			FROM
				chapters AS c
				INNER JOIN titles AS t ON t.id = c.title_id
				INNER JOIN teams ON teams.id = c.team_id
				INNER JOIN user_viewed_chapters AS uvc ON uvc.chapter_id = c.id
			WHERE
				uvc.user_id = ? AND NOT c.hidden
			ORDER BY
				t.id, uvc.created_at DESC
		) AS res
		ORDER BY res.was_read_at DESC
		LIMIT ? OFFSET ?`,
		claims.ID, params.Limit, offset,
	).Scan(&chapters).Error

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(chapters) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено прочитанных вами глав"})
		return
	}

	c.JSON(200, &chapters)
}
