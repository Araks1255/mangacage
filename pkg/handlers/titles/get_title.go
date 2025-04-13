package titles

import (
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func (h handler) GetTitle(c *gin.Context) {
	titleName := c.Param("title")

	var titleID uint
	h.DB.Raw("SELECT id FROM titles WHERE lower(name) = lower(?) AND NOT on_moderation", titleName).Scan(&titleID)
	if titleID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден"})
		return
	}

	var title struct {
		ID          uint
		Name        string
		Description string
		Author      string
		Team        string
		Genres      pq.StringArray `gorm:"type:text[]"`
	}

	h.DB.Raw(
		`SELECT t.id, t.name, t.description, authors.name AS author, teams.name AS team,
		(
		SELECT ARRAY(
		SELECT genres.name FROM genres
		INNER JOIN title_genres ON genres.id = title_genres.genre_id
		INNER JOIN titles ON title_genres.title_id = titles.id
		WHERE titles.id = t.id) AS genres
		)
		FROM titles AS t
		INNER JOIN authors ON authors.id = t.author_id
		INNER JOIN teams ON teams.id = t.team_id
		WHERE NOT t.on_moderation
		AND t.id = ?`, titleID).Scan(&title)

	c.JSON(200, title)
}
