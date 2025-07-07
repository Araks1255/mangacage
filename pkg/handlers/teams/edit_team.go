package teams

import (
	"errors"
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	"github.com/Araks1255/mangacage/pkg/handlers/helpers"
	"github.com/Araks1255/mangacage/pkg/handlers/helpers/teams"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func (h handler) EditTeam(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var requestBody *dto.EditTeamDTO

	team, err := mapEditTeamBodyToTeamOnModeration(requestBody, c.ShouldBindWith, claims.ID)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	if team.Name != nil {
		var teamWithTheSameNameExists bool

		err := tx.Raw("SELECT EXISTS(SELECT 1 FROM teams WHERE lower(name) = lower(?))", team.Name).Scan(&teamWithTheSameNameExists).Error
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}

		if teamWithTheSameNameExists {
			c.AbortWithStatusJSON(409, gin.H{"error": "команда с таким названием уже существует"})
			return
		}
	}

	if err := tx.Raw("SELECT team_id FROM users WHERE id = ?", claims.ID).Scan(&team.ExistingID).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	err = tx.Clauses(helpers.OnIDConflictClause).Create(&team).Error

	if err != nil {
		code, err := parseEditTeamError(err)
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	if requestBody.Cover != nil {
		err := teams.UpsertTeamOnModerationCover(c.Request.Context(), h.TeamsCovers, requestBody.Cover, team.ID, claims.ID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "изменения команды успешно отправлены на модерацию"})
	// Уведомление
}

func mapEditTeamBodyToTeamOnModeration(requestBody *dto.EditTeamDTO, bindFn func(any, binding.Binding) error, userID uint) (*models.TeamOnModeration, error) {
	if err := bindFn(requestBody, binding.FormMultipart); err != nil {
		return nil, err
	}

	if requestBody.Name == nil && requestBody.Description == nil && requestBody.Cover == nil {
		return nil, errors.New("необходим как минимум 1 изменямый параметр")
	}

	if requestBody.Cover != nil && requestBody.Cover.Size > 2<<20 {
		return nil, errors.New("превышен максимальный размер обложки (2мб)")
	}

	res := requestBody.ToTeamOnModeration(userID, 0)

	return &res, nil
}

func parseEditTeamError(err error) (code int, parsedErr error) {
	if dbErrors.IsUniqueViolation(err, constraints.UniqTeamsOnModerationCreatorID) {
		return 409, errors.New("у вас уже есть команда, ожидающая модерации")
	}

	if dbErrors.IsUniqueViolation(err, constraints.UniqTeamsOnModerationName) {
		return 409, errors.New("команда с таким названием уже ожидает модерации")
	}

	return 500, err
}
