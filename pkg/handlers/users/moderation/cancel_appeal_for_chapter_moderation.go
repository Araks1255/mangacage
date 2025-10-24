package moderation

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h handler) CancelAppealForChapterModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	chapterOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id главы на модерации"})
		return
	}

	pagesDirPath, code, err := deleteChapterOnModeration(h.DB, uint(chapterOnModerationID), claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": "заявка на модерацию главы успешно отменена"})

	if pagesDirPath != nil {
		if err := deleteChapterOnModerationPages(*pagesDirPath); err != nil {
			log.Println(err)
		}
	}
}

func deleteChapterOnModeration(db *gorm.DB, id, creatorID uint) (pagesDirPath *string, code int, err error) {
	var res struct {
		ID           *uint
		PagesDirPath *string
	}

	err = db.Raw(
		`DELETE FROM
			chapters_on_moderation
		WHERE
			id = ? AND creator_id = ?
		RETURNING
			id, pages_dir_path`,
		id, creatorID,
	).Scan(&res).Error

	if err != nil {
		return nil, 500, err
	}

	if res.ID == nil {
		return nil, 404, errors.New("глава на модерации не найдена среди созданных вами")
	}

	return res.PagesDirPath, 0, nil
}

func deleteChapterOnModerationPages(path string) error {
	if err := os.RemoveAll(path); err != nil {
		return fmt.Errorf("ошибка при удалении страниц главы на модерации\nпуть: %s\nошибка: %s", path, err.Error())
	}
	return nil
}
