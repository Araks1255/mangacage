package chapters

import (
	"context"
	"database/sql"
	"log"
	"slices"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
)

func (h handler) EditChapter(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	desiredChapterID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id главы должен быть числом"})
		return
	}

	var requestBody struct {
		Name            string `json:"name"`
		Description     string `json:"description"`
		DesiredVolumeID uint   `json:"volumeId"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	if requestBody.Name == "" && requestBody.Description == "" && requestBody.DesiredVolumeID == 0 {
		c.AbortWithStatusJSON(400, gin.H{"error": "необходим хотя-бы один изменяемый параметр"})
		return
	}

	var userRoles []string
	h.DB.Raw(
		`SELECT roles.name FROM roles
		INNER JOIN user_roles ON roles.id = user_roles.role_id
		WHERE user_roles.user_id = ?`, claims.ID,
	).Scan(&userRoles)

	if !slices.Contains(userRoles, "team_leader") && !slices.Contains(userRoles, "ex_team_leader") {
		c.AbortWithStatusJSON(403, gin.H{"error": "у вас недостаточно прав для редактирования главы"})
		return
	}

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var titleID, volumeID, existingChapterID uint

	row := tx.Raw(
		`SELECT
			t.id, v.id, c.id
		FROM
			chapters AS c
			INNER JOIN volumes AS v ON v.id = c.volume_id
			INNER JOIN titles AS t ON t.id = v.title_id
		WHERE
			c.id = ? AND t.team_id = (SELECT team_id FROM users WHERE id = ?)`,
		desiredChapterID, claims.ID,
	).Row()

	if err = row.Scan(&titleID, &volumeID, &existingChapterID); err != nil {
		log.Println(err)
	}
	if existingChapterID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "глава не найдена среди переводимых вашей командой тайтлов"})
		return
	}

	if requestBody.Name != "" {
		var chapterOnModerationWithTheSameNameID, chapterWithTheSameNameID sql.NullInt64

		row = tx.Raw(
			`SELECT
				(SELECT id FROM chapters_on_moderation WHERE volume_id = ? AND lower(name) = lower(?) LIMIT 1),
				(SELECT id FROM chapters WHERE volume_id = ? AND lower(name) = lower(?) LIMIT 1)`,
			volumeID, requestBody.Name, volumeID, requestBody.Name,
		).Row()

		if err = row.Scan(&chapterOnModerationWithTheSameNameID, &chapterWithTheSameNameID); err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}

		if chapterOnModerationWithTheSameNameID.Valid {
			c.AbortWithStatusJSON(409, gin.H{"error": "глава с таким названием уже ожидает модерации в этом томе"})
			return
		}
		if chapterWithTheSameNameID.Valid {
			c.AbortWithStatusJSON(409, gin.H{"error": "глава с таким названием уже есть в этом томе"})
			return
		}
	}

	editedChapter := models.ChapterOnModeration{
		ExistingID:  sql.NullInt64{Int64: int64(existingChapterID), Valid: true},
		Name:        requestBody.Name,
		Description: requestBody.Description,
		CreatorID:   claims.ID,
	}
	if requestBody.DesiredVolumeID != 0 {
		tx.Raw("SELECT id FROM volumes WHERE id = ? AND title_id = ?", requestBody.DesiredVolumeID, titleID).Scan(&editedChapter.VolumeID)
		if !editedChapter.VolumeID.Valid {
			c.AbortWithStatusJSON(404, gin.H{"error": "том не найден"})
			return
		}
	}

	if result := tx.Exec(
		`INSERT INTO chapters_on_moderation (created_at, name, description, creator_id, volume_id, existing_id)
		VALUES (NOW(), ?, ?, ?, ?, ?)
		ON CONFLICT (existing_id) DO UPDATE
		SET
			updated_at = NOW(),
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			creator_id = EXCLUDED.creator_id,
			volume_id = EXCLUDED.volume_id`,
		editedChapter.Name, editedChapter.Description, editedChapter.CreatorID, editedChapter.VolumeID, editedChapter.ExistingID,
	); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "изменения главы успешно отправлены на модерацию"})

	if _, err := h.NotificationsClient.NotifyAboutChapterOnModeration(context.TODO(), &pb.ChapterOnModeration{ID: uint64(editedChapter.ExistingID.Int64), New: false}); err != nil {
		log.Println(err)
	}
}
