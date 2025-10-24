package moderation

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
)

func (h handler) GetMyChapterOnModerationPage(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	chapterOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id главы на модерации"})
		return
	}

	numberOfPage, err := strconv.Atoi(c.Param("page"))
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный номер страницы"})
		return
	}

	var path *string

	err = h.DB.Raw(
		`SELECT
			p.path
		FROM
			pages AS p
			INNER JOIN chapters_on_moderation AS com ON com.id = p.chapter_on_moderation_id
		WHERE
			com.id = ? AND com.creator_id = ? AND number = ?`,
		chapterOnModerationID, claims.ID, numberOfPage,
	).Scan(&path).Error

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error":err.Error()})
		return
	}

	if path == nil {
		c.AbortWithStatusJSON(404, gin.H{"error":"страница главы на модерации не найдена"})
		return
	}

	c.File(*path)
}
