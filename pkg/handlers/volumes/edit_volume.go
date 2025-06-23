package volumes

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/common/utils"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm/clause"
)

func (h handler) EditVolume(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	volumeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id тома"})
		return
	}

	var requestBody models.VolumeOnModerationDTO
	c.ShouldBindJSON(&requestBody)

	ok, err := utils.HasAnyNonEmptyFields(&requestBody)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if !ok {
		c.AbortWithStatusJSON(400, gin.H{"error": "запрос должен содержать как минимум 1 изменяемый параметр"})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var check struct {
		TitleID                        uint
		DoesVolumeWithTheSameNameExist bool
	}

	if err := tx.Raw(
		`SELECT
			(
				SELECT tt.title_id FROM title_teams AS tt
				INNER JOIN volumes AS v ON v.title_id = tt.title_id
				INNER JOIN users AS u ON u.team_id = tt.team_id
				WHERE v.id = ? AND u.id = ?
			) AS title_id,
			EXISTS(SELECT 1 FROM volumes WHERE title_id = (SELECT title_id FROM volumes WHERE id = ?) AND lower(name) = lower(?)) AS does_volume_with_the_same_name_exist`,
		volumeID, claims.ID, volumeID, requestBody.Name,
	).Scan(&check).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if check.TitleID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "том не найден среди томов тайтлов, переводимых вашей командой"})
		return
	}
	if check.DoesVolumeWithTheSameNameExist {
		c.AbortWithStatusJSON(409, gin.H{"error": "том с таким названием уже существует в тайтле"})
		return
	}

	volumeIDuint := uint(volumeID)
	editedVolume := requestBody.ToVolumeOnModeration(claims.ID, &check.TitleID, &volumeIDuint)

	onConflictClause := clause.OnConflict{
		Columns:   []clause.Column{{Name: "existing_id"}},
		UpdateAll: true,
	}

	err = tx.Clauses(onConflictClause).Create(&editedVolume).Error

	if err != nil {
		if dbErrors.IsUniqueViolation(err, constraints.UniqVolumeTitle) {
			c.AbortWithStatusJSON(409, gin.H{"error": "том с таким названием уже ожидает модерации в этом тайтле"})
		} else {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "изменения тома успешно отправлены на модерацию"})

	if _, err := h.NotificationsClient.NotifyAboutVolumeOnModeration(c.Request.Context(), &pb.VolumeOnModeration{ID: volumeID, New: false}); err != nil {
		log.Println(err)
	}
}
