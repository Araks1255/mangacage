package moderation

import (
	"time"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetSelfEditedChaptersOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var chapters []struct {
		CreatedAt   time.Time
		Name        string
		Description string
		Existing    string
		Volume      string
	}

	h.DB.Raw(
		`SELECT c.created_at, c.name, c.description,
		chapters.name AS existing, volumes.name AS volume
		FROM chapters_on_moderation AS c
		INNER JOIN chapters ON chapters.id = c.existing_id
		INNER JOIN volumes ON volumes.id = c.volume_id
		WHERE c.creator_id = ?`, claims.ID,
	).Scan(&chapters)

	if len(chapters) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено ваших отредактированных глав на модерации"})
		return
	}

	response := make([]map[string]any, len(chapters), len(chapters))
	for i := 0; i < len(chapters); i++ {
		response[i] = make(map[string]any, 5)

		if chapters[i].Name != "" {
			response[i]["newName"] = chapters[i].Name
		}
		if chapters[i].Description != "" {
			response[i]["newDescription"] = chapters[i].Description
		}

		response[i]["existing"] = chapters[i].Existing
		response[i]["volume"] = chapters[i].Volume
		response[i]["createdAt"] = chapters[i].CreatedAt.Format(time.DateTime)
	}

	c.JSON(200, &response)
}
