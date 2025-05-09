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
)

func (h handler) CreateTeam(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	if len(form.Value["name"]) == 0 {
		c.AbortWithStatusJSON(400, gin.H{"error": "в запросе не хватает названия команды"})
		return
	}
	if len(form.File["cover"]) == 0 {
		c.AbortWithStatusJSON(400, gin.H{"error": "в запросе не хватает обложки команды"})
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

	var check struct {
		DoesUserHaveTeam             bool
		DoesTeamWithTheSameNameExist bool
	}

	if err := tx.Raw(
		`SELECT
			EXISTS(SELECT 1 FROM users WHERE id = ? AND team_id IS NOT NULL) AS does_user_have_team,
			EXISTS(SELECT 1 FROM teams WHERE lower(name) = lower(?)) AS does_team_with_the_same_name_exist`,
		claims.ID, name,
	).Scan(&check).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if check.DoesUserHaveTeam {
		c.AbortWithStatusJSON(409, gin.H{"error": "вы уже состоите в команде перевода"})
		return
	}
	if check.DoesTeamWithTheSameNameExist {
		c.AbortWithStatusJSON(409, gin.H{"error": "команда с таким названием уже существует"})
		return
	}

	newTeam := models.TeamOnModeration{
		Name:        sql.NullString{String: name, Valid: true},
		Description: description,
		CreatorID:   claims.ID,
	}

	err = tx.Create(&newTeam).Error

	if err != nil {
		if dbErrors.IsUniqueViolation(err, constraints.UniTeamsOnModerationCreatorID) {
			c.AbortWithStatusJSON(409, gin.H{"error": "у вас уже есть команда, ожидающая модерации"})
			return
		}

		if dbErrors.IsUniqueViolation(err, constraints.UniTeamsOnModerationName) {
			c.AbortWithStatusJSON(409, gin.H{"error": "команда с таким названием уже ожидает модерации"})
			return
		}

		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
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

	if _, err := h.TeamsOnModerationCovers.InsertOne(c.Request.Context(), teamCover); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "команда успешно отправлена на модерацию"})
	// Уведомление
}
