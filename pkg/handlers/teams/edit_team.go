package teams

import (
	"database/sql"
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/common/utils"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h handler) EditTeam(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	if len(form.Value["name"]) == 0 && len(form.Value["description"]) == 0 && len(form.File["cover"]) == 0 {
		c.AbortWithStatusJSON(400, gin.H{"error": "запрос должен содержать хотя-бы один изменяемый параметр"})
		return
	}

	if len(form.File["cover"]) != 0 && form.File["cover"][0].Size > 10<<20 {
		c.AbortWithStatusJSON(400, gin.H{"error": "превышен лимит размера обложки (10мб)"})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var userTeamID uint // Тут uint, потому-что сверху уже была проверка на роль лидера команды, а значит, команда у юзера есть, и делать sql.NullInt64 необязательно
	if err := tx.Raw("SELECT team_id FROM users WHERE id = ?", claims.ID).Scan(&userTeamID).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(form.Value["name"]) != 0 {
		var doesTeamWithTheSameNameExist bool

		if err := tx.Raw("SELECT EXISTS(SELECT 1 FROM teams WHERE lower(name) = lower(?))", form.Value["name"][0]).Scan(&doesTeamWithTheSameNameExist).Error; err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}

		if doesTeamWithTheSameNameExist {
			c.AbortWithStatusJSON(409, gin.H{"error": "команда с таким названием уже существует"})
			return
		}
	}

	editedTeam := models.TeamOnModeration{ // Тут можно было просто на переменных сделать, но со структурой мне побольше нравится
		Description: form.Value["description"][0],
		CreatorID:   claims.ID,
		ExistingID:  sql.NullInt64{Int64: int64(userTeamID), Valid: true},
	}

	if len(form.Value["name"]) != 0 {
		editedTeam.Name = sql.NullString{String: form.Value["name"][0]}
	}

	err = tx.Raw(
		`INSERT INTO teams_on_moderation (created_at, name, description, existing_id, creator_id)
		VALUES (NOW(), ?, ?, ?, ?)
		ON CONFLICT (existing_id) DO UPDATE
		SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			creator_id = EXCLUDED.creator_id,
			updated_at = NOW()
		RETURNING id`,
		editedTeam.Name, editedTeam.Description, editedTeam.ExistingID, editedTeam.CreatorID,
	).Scan(&editedTeam.ID).Error

	if err != nil {
		if dbErrors.IsUniqueViolation(err, constraints.UniTeamsOnModerationName) {
			c.AbortWithStatusJSON(409, gin.H{"error": "команда с таким названием уже ожидает модерации"})
		} else {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}
		return
	}

	if len(form.File["cover"]) != 0 {
		cover, err := utils.ReadMultipartFile(form.File["cover"][0], 10<<20)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}

		filter := bson.M{"team_on_moderation_id": editedTeam.ID}
		update := bson.M{"$set": bson.M{"cover": cover}}
		opts := options.Update().SetUpsert(true)

		if _, err := h.TeamsOnModerationCovers.UpdateOne(c.Request.Context(), filter, update, opts); err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "изменения команды успешно отправлены на модерацию"})
	// Уведомление
}
