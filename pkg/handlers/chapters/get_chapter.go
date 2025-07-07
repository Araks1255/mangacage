package chapters

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
)

func (h handler) GetChapter(c *gin.Context) {
	chapterID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id главы"})
		return
	}

	var chapter dto.ResponseChapterDTO

	if err := h.DB.Raw(
		`SELECT 
			c.*,
			v.name AS volume,
			teams.name AS team,
			t.name AS title, t.id AS title_id
		FROM
			chapters AS c
			INNER JOIN volumes AS v ON v.id = c.volume_id
			INNER JOIN titles AS t ON t.id = v.title_id
			INNER JOIN teams ON teams.id = c.team_id
		WHERE
			c.id = ?`,
		chapterID,
	).Scan(&chapter).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if chapter.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "глава не найдена"})
		return
	}

	c.JSON(200, &chapter)
}
