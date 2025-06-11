package chapters

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

func (h handler) EditChapter(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	chapterID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id главы"})
		return
	}

	var requestBody struct {
		Name        string `json:"name"`
		Description string `json:"description"`
		VolumeID    uint   `json:"volumeId"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	if requestBody.Name == "" && requestBody.Description == "" && requestBody.VolumeID == 0 {
		c.AbortWithStatusJSON(400, gin.H{"error": "необходим хотя-бы один изменяемый параметр"})
		return
	}

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var doesChapterExist bool
	if err := tx.Raw(
		`SELECT EXISTS (
			SELECT 1
			FROM chapters AS c
			INNER JOIN volumes AS v ON v.id = c.volume_id
			INNER JOIN titles AS t ON t.id = v.title_id
			WHERE c.id = ? AND t.team_id = (SELECT team_id FROM users WHERE id = ?)
		)`,
		chapterID, claims.ID,
	).Scan(&doesChapterExist).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if !doesChapterExist {
		c.AbortWithStatusJSON(404, gin.H{"error": "глава не найдена среди глав, переводимых вашей командой"})
		return
	}

	chapterIDuint := uint(chapterID)

	editedChapter := models.ChapterOnModeration{
		CreatorID:   claims.ID,
		Description: requestBody.Description,
		ExistingID:  &chapterIDuint,
	}

	if requestBody.VolumeID != 0 {
		editedChapter.VolumeID = requestBody.VolumeID
	} else {
		if err := tx.Raw("SELECT volume_id FROM chapters WHERE id = ?", chapterID).Scan(&editedChapter.VolumeID).Error; err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	if requestBody.Name != "" {
		var doesChapterWithThisNameAlreadyExist bool

		if err := tx.Raw(
			"SELECT EXISTS(SELECT 1 FROM chapters WHERE lower(name) = lower(?) AND volume_id = ?)",
			requestBody.Name, editedChapter.VolumeID,
		).Scan(&doesChapterWithThisNameAlreadyExist).Error; err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}

		if doesChapterWithThisNameAlreadyExist {
			c.AbortWithStatusJSON(409, gin.H{"error": "глава с таким названием уже существует в этом томе"})
			return
		}

		editedChapter.Name = sql.NullString{String: requestBody.Name, Valid: true}
	}

	err = tx.Exec(
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
	).Error

	if err != nil {
		if dbErrors.IsUniqueViolation(err, constraints.UniqChapterVolume) {
			c.AbortWithStatusJSON(409, gin.H{"error": "глава с таким названием уже ожидает модерации в этом томе"})
			return
		}

		if dbErrors.IsForeignKeyViolation(err, constraints.FkChaptersOnModerationVolume) {
			c.AbortWithStatusJSON(404, gin.H{"error": "том не найден"})
			return
		}

		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "изменения главы успешно отправлены на модерацию"})

	if _, err := h.NotificationsClient.NotifyAboutChapterOnModeration(c.Request.Context(), &pb.ChapterOnModeration{ID: uint64(*editedChapter.ExistingID), New: false}); err != nil {
		log.Println(err)
	}
}
