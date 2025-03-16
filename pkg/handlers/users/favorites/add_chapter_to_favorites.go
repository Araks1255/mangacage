package favorites

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) AddChapterToFavorites(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var requestBody struct {
		Title   string `json:"title" binding:"required"`
		Volume  string `json:"volume" binding:"required"`
		Chapter string `json:"chapter" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	var chapterID uint
	h.DB.Raw(`SELECT chapters.id FROM chapters
		INNER JOIN volumes ON chapters.volume_id = volumes.id
		INNER JOIN titles ON volumes.title_id = titles.id
		WHERE lower(titles.name) = lower(?)
		AND lower(volumes.name) = lower(?)
		AND lower(chapters.name) = lower(?)
		AND not chapters.on_moderation`,
		requestBody.Title, requestBody.Volume, requestBody.Chapter,
	).Scan(&chapterID)

	if chapterID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "глава не найдена"})
		return
	}

	if result := h.DB.Exec("INSERT INTO user_favorite_chapters (user_id, chapter_id) VALUES (?,?)", claims.ID, chapterID); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(201, gin.H{"success": "глава успешно добавлена в избранное"})
}
