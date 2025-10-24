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

func (h handler) CancelAppealForProfileChanges(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	profilePicturePath, code, err := deleteProfileChanges(h.DB, claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": "заявка на модерацию изменений профиля успешно отменена"})

	if profilePicturePath != nil {
		if err := deleteProfileChangesProfilePicture(*profilePicturePath); err != nil {
			log.Println(err)
		}
	}
}

func deleteProfileChanges(db *gorm.DB, existingID uint) (profilePicturePath *string, code int, err error) {
	var res struct {
		ID                 *uint
		ProfilePicturePath *string
	}

	err = db.Raw(
		`DELETE FROM
			users_on_moderation
		WHERE
			existing_id = ?
		RETURNING
			id, profile_picture_path`,
	).Scan(&res).Error

	if err != nil {
		return nil, 500, err
	}

	if res.ID == nil {
		return nil, 404, errors.New("не найдено ваших изменений профиля")
	}

	return res.ProfilePicturePath, 0, nil
}

func deleteProfileChangesProfilePicture(path string) error {
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("ошибка при удалении аватарки изменений профиля\nпуть: %s\nошибка: %s", path, err.Error())
	}
	return nil
}
