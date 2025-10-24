package titles

import (
	"log"
	"strconv"
	"strings"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
)

func (h handler) GetTitle(c *gin.Context) {
	titleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id тайтла"})
		return
	}

	var selects strings.Builder
	args := make([]any, 0, 2)

	selects.WriteString(
		`t.*,
		a.name AS author,
		COALESCE(ARRAY_AGG(DISTINCT g.name)::TEXT[], '{}'::TEXT[]) AS genres,
		COALESCE(ARRAY_AGG(DISTINCT tags.name)::TEXT[], '{}'::TEXT[]) AS tags,
		ARRAY(SELECT DISTINCT volume FROM chapters WHERE title_id = t.id AND NOT hidden ORDER BY volume DESC) AS volumes,
		ARRAY(SELECT DISTINCT team_id FROM title_teams WHERE title_id = t.id) AS teams_ids`,
	)

	claims, ok := c.Get("claims")
	if ok {
		selects.WriteString(
			`,COUNT(DISTINCT uvc.*) AS quantity_of_viewed_chapters,
			tr.rate AS user_rate,
			EXISTS(SELECT 1 FROM user_titles_subscribed_to WHERE user_id = ? AND title_id = t.id) AS my_subscription,
			EXISTS(SELECT 1 FROM user_favorite_titles WHERE user_id = ? AND title_id = t.id) AS favorited_by_me`,
		)
		args = append(args, claims.(*auth.Claims).ID, claims.(*auth.Claims).ID)
	}

	query := h.DB.Table("titles AS t").
		Select(selects.String(), args...).
		Joins("INNER JOIN authors AS a ON t.author_id = a.id").
		Joins("LEFT JOIN title_genres AS tg ON t.id = tg.title_id").
		Joins("INNER JOIN genres AS g ON tg.genre_id = g.id").
		Joins("LEFT JOIN title_tags AS tt ON t.id = tt.title_id").
		Joins("INNER JOIN tags ON tt.tag_id = tags.id").
		Group("t.id, a.id").
		Where("NOT t.hidden").
		Where("t.id = ?", titleID)

	if ok {
		query = query.Joins("LEFT JOIN chapters AS c ON t.id = c.title_id AND NOT c.hidden").
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
