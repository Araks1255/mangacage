package volumes

import (
	"database/sql"
	"log"
	"slices"
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

	var userRoles []string
	h.DB.Raw(
		`SELECT r.name FROM roles AS r
		INNER JOIN user_roles AS ur ON ur.role_id = r.id
		WHERE ur.user_id = ? AND r.type = 'team'`,
		claims.ID,
	).Scan(&userRoles)

	if !slices.Contains(userRoles, "team_leader") && !slices.Contains(userRoles, "ex_team_leader") {
		c.AbortWithStatusJSON(403, gin.H{"error": "у вас недостаточно прав для редактирования тома"})
		return
	}

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var titleID sql.NullInt64

	if err := tx.Raw(
		`SELECT t.id
		FROM
			titles AS t
			INNER JOIN volumes AS v ON t.id = v.title_id
			INNER JOIN users AS u ON t.team_id = u.team_id
		WHERE
			v.id = ? AND u.id = ?`,
		volumeID, claims.ID,
	).Scan(&titleID).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if !titleID.Valid {
		c.AbortWithStatusJSON(404, gin.H{"error": "том не найден среди тайтлов, переводимых вашей командой"})
		return
	}

	editedVolume := models.VolumeOnModeration{
		ExistingID:  sql.NullInt64{Int64: int64(volumeID), Valid: true},
		Description: requestBody.Description,
		TitleID:     uint(titleID.Int64),
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
