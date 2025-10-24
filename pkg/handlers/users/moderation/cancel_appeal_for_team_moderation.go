package moderation

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h handler) CancelAppealForTeamModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	coverPath, code, err := deleteTeamOnModeration(h.DB, claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": "заявка на модерацию команды успешно отменена"})

	if coverPath != nil {
		if err := deleteTeamOnModerationCover(*coverPath); err != nil {
			log.Println(err)
		}
	}
}

func deleteTeamOnModeration(db *gorm.DB, creatorID uint) (coverPath *string, code int, err error) {
	var res struct {
		ID        *uint
		CoverPath *string
	}

	err = db.Raw(
		`DELETE FROM
			teams_on_moderation
		WHERE
			creator_id = ?
		RETURNING
			id, cover_path`,
	).Scan(&res).Error

	if err != nil {
		return nil, 500, err
	}

	if res.ID == nil {
		return nil, 404, errors.New("не найдено вашей команды на модерации")
	}

	return res.CoverPath, 0, nil
}

func deleteTeamOnModerationCover(path string) error {
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("не удалось удалить обложку команды на модерации\nпуть: %s\nошибка: %s", path, err.Error())
	}
	return nil
}
