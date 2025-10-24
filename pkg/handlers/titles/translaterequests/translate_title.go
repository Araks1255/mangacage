package translaterequests

import (
	"errors"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/handlers/helpers/titles"
	pb "github.com/Araks1255/mangacage_protos/gen/site_notifications"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h handler) TranslateTitle(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	titleID, isTitleTranslating, code, err := parseTranslateTitleParams(h.DB, claims.ID, c.Param)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	if !isTitleTranslating {
		code, err := translateTitle(tx, titleID, claims.ID)

		if err != nil {
			if code == 500 {
				log.Println(err)
			}
			c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}

		c.JSON(201, gin.H{"success": "теперь ваша команда переводит этот тайтл"})
	} else {
		var requestBody struct {
			Message *string `json:"message"`
		}

		c.ShouldBindJSON(&requestBody) // Тело запроса опциональное

		err := tx.Exec(
			`INSERT INTO titles_translate_requests (title_id, team_id, created_at, message)
			VALUES (?, (SELECT team_id FROM users WHERE id = ?), NOW(), ?)`,
			titleID, claims.ID, requestBody.Message,
		).Error

		if err != nil {
			code, err := titles.ParseTitleTranslateRequestInsertError(err)
			if code == 500 {
				log.Println(err)
			}
			c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}

		c.JSON(201, gin.H{"success": "заявка на перевод тайтла успешно отправлена на модерацию"})

		var message string
		if requestBody.Message != nil {
			message = *requestBody.Message
		}

		if _, err := h.NotificationsClient.NotifyAboutTitleTranslateRequest(
			c.Request.Context(), &pb.TitleTranslateRequest{
				TitleID:  uint64(titleID),
				SenderID: uint64(claims.ID),
				Message:  message,
			},
		); err != nil {
			log.Println(err)
		}
	}

	tx.Commit()
}

func parseTranslateTitleParams(db *gorm.DB, userID uint, paramFn func(string) string) (titleID uint, isTitleTranslating bool, code int, err error) {
	titleIDuint64, err := strconv.ParseUint(paramFn("id"), 10, 64)
	if err != nil {
		return 0, false, 400, errors.New("указан невалидный id тайтла")
	}

	err = db.Raw(
		"SELECT EXISTS(SELECT 1 FROM title_teams WHERE title_id = ? AND team_id != (SELECT team_id FROM users WHERE id = ?))",
		titleIDuint64, userID,
	).Scan(&isTitleTranslating).Error
	if err != nil {
		return 0, false, 500, err
	}

	return uint(titleIDuint64), isTitleTranslating, 0, nil
}

func translateTitle(db *gorm.DB, titleID, userID uint) (code int, err error) {
	err = db.Exec("INSERT INTO title_teams (title_id, team_id) VALUES (?, (SELECT team_id FROM users WHERE id = ?))", titleID, userID).Error
	if err != nil {
		return titles.ParseTitleTeamInsertError(err)
	}

	result := db.Exec("UPDATE titles SET translating_status = 'ongoing' WHERE id = ?", titleID)
	if result.Error != nil {
		return 500, result.Error
	}

	if result.RowsAffected == 0 {
		return 500, errors.New("не удалось обновить статус перевода тайтла")
	}

	return 0, nil
}
