package titles

import (
	"fmt"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
)

func (h handler) GetTitle(c *gin.Context) {
	desiredTitleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id тайтла"})
		return
	}

	selects := []string{
		"t.*",
		"a.name AS author",
		"ARRAY_AGG(DISTINCT g.name)::TEXT[] AS genres",
		"ARRAY_AGG(DISTINCT tags.name)::TEXT[] AS tags",
		"ROUND(t.sum_of_rates::numeric / NULLIF(t.number_of_rates, 0), 1) AS rate",
	}

	claims, ok := c.Get("claims")
	if ok {
		selects = append(
			selects,
			"COUNT(DISTINCT uvc.*) AS quantity_of_viewed_chapters",
			"COUNT(DISTINCT c.id) AS qunatity_of_chapters",
			"tr.rate AS user_rate",
			fmt.Sprintf(
				`EXISTS(
					SELECT 1 FROM users AS u
					INNER JOIN user_roles AS ur ON ur.user_id = u.id
					INNER JOIN roles AS r ON r.id = ur.role_id
					INNER JOIN title_teams ON title_teams.team_id = u.team_id
					WHERE title_teams.title_id = t.id AND u.id = %d AND r.name IN ('team_leader', 'ex_team_leader')
				) AS can_edit`, claims.(*auth.Claims).ID,
			),
		)
	}

	query := h.DB.Table("titles AS t").Select(selects).
		Joins("INNER JOIN authors AS a ON t.author_id = a.id").
		Joins("INNER JOIN title_genres AS tg ON tg.title_id = t.id").
		Joins("INNER JOIN genres AS g ON tg.genre_id = g.id").
		Joins("INNER JOIN title_tags AS tt ON t.id = tt.title_id").
		Joins("INNER JOIN tags ON tt.tag_id = tags.id").
		Group("t.id, a.id").
		Where("t.id = ?", desiredTitleID)

	if ok {
		query = query.Joins("LEFT JOIN volumes AS v ON t.id = v.title_id").
			Joins("LEFT JOIN chapters AS c ON v.id = c.volume_id").
			Joins("LEFT JOIN user_viewed_chapters AS uvc ON c.id = uvc.chapter_id AND uvc.user_id = ?", claims.(*auth.Claims).ID).
			Joins("LEFT JOIN title_rates AS tr ON t.id = tr.title_id AND tr.user_id = ?", claims.(*auth.Claims).ID).
			Group("tr.rate")
	}

	var title dto.ResponseTitleDTO

	if err := query.Scan(&title).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if title.ID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден"})
		return
	}

	c.JSON(200, &title)
}
