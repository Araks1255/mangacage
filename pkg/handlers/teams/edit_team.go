package teams

import (
	"errors"
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/Araks1255/mangacage/pkg/common/utils"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	"github.com/Araks1255/mangacage/pkg/handlers/helpers"
	"github.com/Araks1255/mangacage/pkg/handlers/helpers/teams"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

func (h handler) EditTeam(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var requestBody dto.EditTeamDTO

	if err := c.ShouldBindWith(&requestBody, binding.FormMultipart); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	code, err := checkEditTeamConflicts(h.DB, requestBody, claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	existingID, err := getUserTeamID(h.DB, claims.ID)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	team := requestBody.ToTeamOnModeration(claims.ID, existingID)

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	err = tx.Clauses(helpers.OnExistingIDConflictClause).Create(&team).Error

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

func checkEditTeamConflicts(db *gorm.DB, requestBody dto.EditTeamDTO, userID uint) (code int, err error) {
	ok, err := utils.HasAnyNonEmptyFields(requestBody)
	if err != nil {
		return 500, err
	}

	if !ok {
		return 400, errors.New("необходим как минимум 1 изменяемый параметр")
	}

	if requestBody.Cover != nil && requestBody.Cover.Size > 2<<20 {
		return 400, errors.New("превышен максимальный размер обложки (2мб)")
	}

	if requestBody.Name != nil {
		var teamExists bool

		err = db.Raw("SELECT EXISTS(SELECT 1 FROM teams WHERE lower(name) = lower(?))", requestBody.Name).Scan(&teamExists).Error
		if err != nil {
			return 500, err
		}

		if teamExists {
			return 409, errors.New("команда с таким названием уже ожидает модерации")
		}
	}

	return 0, nil
}

func getUserTeamID(db *gorm.DB, userID uint) (uint, error) {
	var res *uint

	if err := db.Raw("SELECT team_id FROM users WHERE id = ?", userID).Scan(&res).Error; err != nil {
		return 0, err
	}

	if res == nil {
		return 0, errors.New("ваша команда не найдена")
	}

	return *res, nil
}

func parseEditTeamError(err error) (code int, parsedErr error) {
	if dbErrors.IsUniqueViolation(err, constraints.UniTeamsOnModerationCreatorID) {
		return 409, errors.New("у вас уже есть команда, ожидающая модерации")
	}

	if dbErrors.IsUniqueViolation(err, constraints.UniqTeamOnModerationName) {
		return 409, errors.New("команда с таким названием уже ожидает модерации")
	}

	return 500, err
}
