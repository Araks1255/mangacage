package moderation

import (
	"time"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) GetSelfEditedVolumesOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var volumes []struct {
		CreatedAt   time.Time
		Name        string
		Description string
		Existing    string
		Title       string
	}

	h.DB.Raw(
		`SELECT v.created_at, v.name, v.description,
		volumes.name AS existing, titles.name AS title
		FROM volumes_on_moderation AS v
		INNER JOIN volumes ON v.existing_id = volumes.id
		INNER JOIN titles ON v.title_id = titles.id
		WHERE v.creator_id = ?`, claims.ID,
	).Scan(&volumes)

	if len(volumes) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено ваших отредактированных томов на модерации"})
		return
	}

	response := make([]map[string]any, len(volumes), len(volumes))
	for i := 0; i < len(volumes); i++ {
		response[i] = make(map[string]any, 5)

		if volumes[i].Name != "" {
			response[i]["newName"] = volumes[i].Name
		}
		if volumes[i].Description != "" {
			response[i]["newDescription"] = volumes[i].Description
		}

		response[i]["existing"] = volumes[i].Existing
		response[i]["title"] = volumes[i].Title
		response[i]["createdAt"] = volumes[i].CreatedAt.Format(time.DateTime)
	}

	c.JSON(200, &response)
}
