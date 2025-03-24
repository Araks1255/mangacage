package moderation

import (
	"time"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

func (h handler) GetSelfEditedTitlesOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var editedTitles []struct {
		CreatedAt   time.Time
		Name        string
		Description string
		Existing    string
		Author      string
		Genres      pq.StringArray `gorm:"type:TEXT[]"`
	}

	h.DB.Raw(
		`SELECT t.created_at, t.name, t.description, titles.name AS existing, authors.name AS author, t.genres
		FROM titles_on_moderation AS t
		INNER JOIN titles ON titles.id = t.existing_id
		LEFT JOIN authors ON authors.id = t.author_id
		WHERE t.creator_id = ?`, claims.ID,
	).Scan(&editedTitles)

	if len(editedTitles) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено ваших отредактированных тайтлов на модерации"})
		return
	}

	response := make([]map[string]any, len(editedTitles), len(editedTitles))
	for i := 0; i < len(editedTitles); i++ {
		response[i] = make(map[string]any, 6)

		if editedTitles[i].Name != "" {
			response[i]["newName"] = editedTitles[i].Name
		}
		if editedTitles[i].Description != "" {
			response[i]["newDescription"] = editedTitles[i].Description
		}
		if editedTitles[i].Author != "" {
			response[i]["newAuthor"] = editedTitles[i].Author
		}

		if len(editedTitles[i].Genres) != 0 {
			response[i]["newGenres"] = editedTitles[i].Genres
		}

		response[i]["createdAt"] = editedTitles[i].CreatedAt.Format(time.DateTime)
		response[i]["existing"] = editedTitles[i].Existing
	}

	c.JSON(200, &response)
}
