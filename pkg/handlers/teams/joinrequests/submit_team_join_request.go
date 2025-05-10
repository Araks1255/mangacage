package joinrequests

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) SubmitTeamJoinRequest(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	desiredTeamID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id команды"})
		return
	}

	var requestBody struct {
		IntroductoryMessage string `json:"introductoryMessage"`
		RoleID              uint   `json:"roleId"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(404, gin.H{"error": err.Error()})
		return
	}

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	if requestBody.RoleID != 0 {
		var check struct {
			DoesUserHaveTeam bool
			DoesRoleExist    bool
		}

		if err := tx.Raw(
			`SELECT
				EXISTS(SELECT 1 FROM users WHERE id = ? AND team_id IS NOT NULL) AS does_user_have_team,
				EXISTS(SELECT 1 FROM roles WHERE id = ? AND type = 'team') AS does_role_exist`,
			claims.ID, requestBody.RoleID,
		).Scan(&check).Error; err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}

		if check.DoesUserHaveTeam {
			c.AbortWithStatusJSON(409, gin.H{"error": "вы уже состоите в команде перевода"})
			return
		}
		if !check.DoesRoleExist {
			c.AbortWithStatusJSON(404, gin.H{"error": "роль не найдена"})
			return
		}
	} else {
		var doesUserHaveTeam bool

		if err := tx.Raw("SELECT EXISTS(SELECT 1 FROM users WHERE id = ? AND team_id IS NOT NULL)", claims.ID).Scan(&doesUserHaveTeam).Error; err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}

		if doesUserHaveTeam {
			c.AbortWithStatusJSON(409, gin.H{"error": "вы уже состоите в команде перевода"})
			return
		}
	}

	joinRequest := models.TeamJoinRequest{
		CandidateID:         claims.ID,
		TeamID:              uint(desiredTeamID),
		IntroductoryMessage: requestBody.IntroductoryMessage,
	}
	if requestBody.RoleID != 0 {
		joinRequest.RoleID = sql.NullInt64{Int64: int64(requestBody.RoleID), Valid: true}
	}

	err = tx.Create(&joinRequest).Error

	if err != nil {
		if dbErrors.IsUniqueViolation(err, "uniq_team_join_request") {
			c.AbortWithStatusJSON(409, gin.H{"error": "вы уже оставили заявку на вступление в эту команду"})
			return
		}

		if dbErrors.IsForeignKeyViolation(err, "fk_team_join_requests_team") {
			c.AbortWithStatusJSON(404, gin.H{"error": "команда не найдена"})
			return
		}

		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "заявка на вступление в команду успешно отправлена"})
	// Уведомление лидеру
}
