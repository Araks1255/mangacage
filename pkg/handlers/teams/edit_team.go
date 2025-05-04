package teams

import (
	"context"
	"database/sql"
	"log"
	"slices"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/common/utils"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h handler) EditTeam(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var userTeamID uint
	h.DB.Raw("SELECT team_id FROM users WHERE id = ?", claims.ID).Scan(&userTeamID)
	if userTeamID == 0 {
		c.AbortWithStatusJSON(409, gin.H{"error": "вы не состоите в команде перевода"})
		return
	}

	var userRoles []string
	h.DB.Raw(
		`SELECT roles.name FROM roles
		INNER JOIN user_roles ON roles.id = user_roles.role_id
		WHERE user_roles.user_id = ?`, claims.ID,
	).Scan(&userRoles)

	if !slices.Contains(userRoles, "team_leader") && !slices.Contains(userRoles, "ex_team_leader") && !slices.Contains(userRoles, "moder") && !slices.Contains(userRoles, "admin") {
		c.AbortWithStatusJSON(403, gin.H{"error": "вы не являетесь владельцем команды перевода"})
		return
	}

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

	if len(form.Value["name"]) != 0 {
		var existingTeamOnModerationID, existingTeamID sql.NullInt64

		row := tx.Raw(
			`SELECT
				(SELECT id FROM teams_on_moderation WHERE lower(name) = lower(?) LIMIT 1),
				(SELECT id FROM teams WHERE lower(name) = lower(?) LIMIT 1)`,
			form.Value["name"][0], form.Value["name"][0],
		).Row()

		if err = row.Scan(&existingTeamOnModerationID, &existingTeamID); err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}

		if existingTeamOnModerationID.Valid {
			c.AbortWithStatusJSON(409, gin.H{"error": "команда с таким названием уже ожидает модерации"})
			return
		}
		if existingTeamID.Valid {
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

	if result := tx.Raw(
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
	).Scan(&editedTeam.ID); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
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

		if _, err := h.TeamsOnModerationCovers.UpdateOne(context.Background(), filter, update, opts); err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "изменения команды успешно отправлены на модерацию"})
}
