package teams

import (
	"errors"
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	"github.com/Araks1255/mangacage/pkg/handlers/helpers"
	"github.com/Araks1255/mangacage/pkg/handlers/helpers/teams"
	"github.com/Araks1255/mangacage_protos/gen/enums"
	pb "github.com/Araks1255/mangacage_protos/gen/site_notifications"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

func (h handler) CreateTeam(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var requestBody dto.CreateTeamDTO

	if err := c.ShouldBindWith(&requestBody, binding.FormMultipart); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	code, err := checkCreateTeamConflicts(h.DB, requestBody, claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	newTeam := requestBody.ToTeamOnModeration(claims.ID)

	err = helpers.UpsertEntityOnModeration(tx, newTeam, newTeam.ID)

	if err != nil {
		code, err := parseCreateTeamError(err)
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	if code, err := teams.CreateTeamOnModerationCover(tx, h.PathToMediaDir, newTeam.ID, requestBody.Cover); err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "команда успешно отправлена на модерацию"})

	if _, err := h.NotificationsClient.NotifyAboutNewModerationRequest(
		c.Request.Context(),
		&pb.ModerationRequest{
			EntityOnModeration: enums.EntityOnModeration_ENTITY_ON_MODERATION_TEAM,
			ID:                 uint64(newTeam.ID),
		},
	); err != nil {
		log.Println(err)
	}
}

func checkCreateTeamConflicts(db *gorm.DB, requestBody dto.CreateTeamDTO, userID uint) (code int, err error) {
	if requestBody.Cover.Size > 2<<20 {
		return 400, errors.New("превышен максимальный размер обложки (2мб)")
	}

	if requestBody.ID != nil {
		isOwner, err := helpers.CheckEntityOnModerationOwnership(db, "teams", *requestBody.ID, userID)
		if err != nil {
			return 500, err
		}

		if !isOwner {
			return 403, errors.New("редактировать заявку на модерацию может только её создатель")
		}
	}

	var check struct {
		DoesUserHaveTeam             bool
		DoesTeamWithTheSameNameExist bool
	}

	err = db.Raw(
		`SELECT
			EXISTS(SELECT 1 FROM users WHERE id = ? AND team_id IS NOT NULL) AS does_user_have_team,
			EXISTS(SELECT 1 FROM teams WHERE lower(name) = lower(?)) AS does_team_with_the_same_name_exist`,
		userID, requestBody.Name,
	).Scan(&check).Error

	if err != nil {
		return 500, err
	}

	if check.DoesUserHaveTeam {
		return 409, errors.New("вы уже состоите в команде перевода")
	}

	if check.DoesTeamWithTheSameNameExist {
		return 409, errors.New("команда с таким названием уже существует")
	}

	return 0, nil
}

func parseCreateTeamError(err error) (code int, parsedErr error) {
	if dbErrors.IsUniqueViolation(err, constraints.UniTeamsOnModerationCreatorID) {
		return 409, errors.New("у вас уже есть команда, ожидающая модерации")
	}

	if dbErrors.IsUniqueViolation(err, constraints.UniqTeamOnModerationName) {
		return 409, errors.New("команда с таким названием уже ожидает модерации")
	}

	return 500, err
}
