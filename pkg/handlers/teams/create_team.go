package teams

import (
	"context"
	"database/sql"
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/common/utils"
	"github.com/gin-gonic/gin"
)

func (h handler) CreateTeam(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	if len(form.Value["name"]) == 0 || len(form.File["cover"]) == 0 {
		c.AbortWithStatusJSON(400, gin.H{"error": "в запросе нет имени команды или её обложки"})
		return
	}

	name := form.Value["name"][0]

	var description string
	if len(form.Value["description"]) != 0 {
		description = form.Value["description"][0]
	}

	coverFileHeader := form.File["cover"][0]
	if coverFileHeader.Size > 10<<20 {
		c.AbortWithStatusJSON(400, gin.H{"error": "превышен лимит размера обложки (10мб)"})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var userTeamID uint
	tx.Raw("SELECT team_id FROM users WHERE id = ?", claims.ID).Scan(&userTeamID)
	if userTeamID != 0 {
		c.AbortWithStatusJSON(409, gin.H{"error": "у вас уже есть команда перевода"})
		return
	}

	var existing struct {
		UserTeamOnModerationID uint
		TeamOnModerationID     uint
		TeamID                 uint
	}

	tx.Raw(
		`SELECT
			(SELECT id FROM teams_on_moderation WHERE creator_id = ? LIMIT 1) AS user_team_on_moderation_id,
			(SELECT id FROM teams_on_moderation WHERE lower(name) = lower(?) LIMIT 1) AS team_on_moderation_id,
			(SELECT id FROM teams WHERE lower(name) = lower(?) LIMIT 1) AS team_id`,
		claims.ID, name, name,
	).Scan(&existing)

	if existing.UserTeamOnModerationID != 0 {
		c.AbortWithStatusJSON(409, gin.H{"error": "у вас уже есть команда, ожидающая модерации"})
		return
	}
	if existing.TeamOnModerationID != 0 {
		c.AbortWithStatusJSON(409, gin.H{"error": "команда с таким названием уже ожидает модерации"})
		return
	}
	if existing.TeamID != 0 {
		c.AbortWithStatusJSON(409, gin.H{"error": "команда с таким названием уже существует"})
		return
	}

	newTeam := models.TeamOnModeration{
		Name:        sql.NullString{String: name, Valid: true},
		Description: description,
		CreatorID:   claims.ID,
	}

	if result := tx.Create(&newTeam); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	var teamCover struct {
		TeamOnModerationID uint   `bson:"team_on_moderation_id"`
		Cover              []byte `bson:"cover"`
	}

	teamCover.TeamOnModerationID = newTeam.ID
	teamCover.Cover, err = utils.ReadMultipartFile(coverFileHeader, 10<<20)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if _, err := h.TeamsOnModerationCovers.InsertOne(context.Background(), teamCover); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "команда успешно отправлена на модерацию"})
	// Уведомление модерам
}
