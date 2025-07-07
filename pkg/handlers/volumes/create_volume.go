package volumes

import (
	"errors"
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	"github.com/Araks1255/mangacage/pkg/handlers/helpers"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h handler) CreateVolume(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var requestBody dto.CreateVolumeDTO

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	if requestBody.TitleID != nil && requestBody.TitleOnModerationID != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "должен быть заполнен только один id тайтла"})
		return
	}

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	code, err := checkCreateVolumeConflicts(tx, &requestBody, claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	volume := requestBody.ToVolumeOnModeration(claims.ID)

	err = tx.Clauses(helpers.OnIDConflictClause).Create(&volume).Error

	if err != nil {
		if dbErrors.IsUniqueViolation(err, constraints.UniqVolumeOnModerationTitle) {
			c.AbortWithStatusJSON(409, gin.H{"error": "том с таким названием уже ожидает модерации в этом тайтле"})
			return
		}

		if dbErrors.IsUniqueViolation(err, constraints.UniqVolumeOnModerationTitleOnModeration) {
			c.AbortWithStatusJSON(409, gin.H{"error": "том с таким названием уже ожидает модерации в этом тайтле на модерации"})
			return
		}

		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "том успешно отправлен на модерацию"})

	if _, err := h.NotificationsClient.NotifyAboutVolumeOnModeration(c.Request.Context(), &pb.VolumeOnModeration{ID: uint64(volume.ID), New: true}); err != nil {
		log.Println(err)
	}
}

func checkCreateVolumeConflicts(db *gorm.DB, volume *dto.CreateVolumeDTO, userID uint) (code int, err error) {
	if volume.TitleID != nil {
		var check struct {
			UserTeamID                  *uint
			VolumeWithTheSameNameExists bool
		}

		err = db.Raw(
			`SELECT
				(SELECT team_id FROM title_teams WHERE title_id = ? AND team_id = (SELECT team_id FROM users WHERE id = ?)) AS user_team_id,
				EXISTS(SELECT 1 FROM volumes WHERE lower(name) = lower(?)) AS does_volume_with_the_same_name_exist`,
			volume.TitleID, userID, volume.Name,
		).Scan(&check).Error

		if err != nil {
			return 500, err
		}

		if check.UserTeamID == nil {
			return 404, errors.New("тайтл не найден среди переводимых вашей командой")
		}

		if check.VolumeWithTheSameNameExists {
			return 409, errors.New("том с таким названием уже существует")
		}

		volume.TeamID = check.UserTeamID
	}

	if volume.TitleOnModerationID != nil {
		var userTeamID *uint

		err = db.Raw(
			`SELECT u.team_id FROM users AS u
			INNER JOIN titles_on_moderation AS tom ON tom.creator_id = u.id
			WHERE tom.id = ? AND u.id = ?`,
			volume.TitleOnModerationID, userID,
		).Scan(&userTeamID).Error

		if err != nil {
			return 500, err
		}

		if userTeamID == nil {
			return 404, errors.New("тайтл на модерации не найден среди ваших заявок")
		}

		volume.TeamID = userTeamID
	}

	return 0, nil
}
