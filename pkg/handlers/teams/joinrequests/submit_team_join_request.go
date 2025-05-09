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

	var userTeamID sql.NullInt64

	if err := tx.Raw("SELECT team_id FROM users WHERE id = ?", claims.ID).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if userTeamID.Valid {
		c.AbortWithStatusJSON(409, gin.H{"error": "вы уже состоите в команде перевода"})
		return
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
		if dbErrors.IsForeignKeyViolation(err, "fk_team_join_requests_role") {
			c.AbortWithStatusJSON(404, gin.H{"error": "роль не найдена"})
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
