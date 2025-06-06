package volumes

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
)

func (h handler) EditVolume(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	volumeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id тома"})
		return
	}

	var requestBody struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	if requestBody.Name == "" && requestBody.Description == "" {
		c.AbortWithStatusJSON(400, gin.H{"error": "необходим хотя-бы один изменяемый параметр"})
		return
	}

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var check struct {
		TitleID                        uint
		DoesVolumeWithTheSameNameExist bool
	}

	if err := tx.Raw(
		`SELECT
			(SELECT id FROM titles WHERE id = (SELECT title_id FROM volumes WHERE id = ?) AND team_id = (SELECT team_id FROM users WHERE id = ?)) AS title_id,
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

	editedVolume := models.VolumeOnModeration{
		ExistingID:  &volumeIDuint,
		Description: requestBody.Description,
		TitleID:     check.TitleID,
		CreatorID:   claims.ID,
	}
	if requestBody.Name != "" {
		editedVolume.Name = sql.NullString{String: requestBody.Name, Valid: true}
	}

	err = tx.Exec(
		`INSERT INTO volumes_on_moderation (created_at, name, description, existing_id, title_id, creator_id)
		VALUES (NOW(), ?, ?, ?, ?, ?)
		ON CONFLICT (existing_id) DO UPDATE
		SET
			updated_at = EXCLUDED.created_at,
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			creator_id = EXCLUDED.creator_id
		RETURNING id`,
		editedVolume.Name, editedVolume.Description, editedVolume.ExistingID, editedVolume.TitleID, editedVolume.CreatorID,
	).Error

	if err != nil {
		if dbErrors.IsUniqueViolation(err, constraints.UniVolumeTitle) {
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
