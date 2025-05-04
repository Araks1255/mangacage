package joinrequests

import (
	"database/sql"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) SubmitTeamJoinRequest(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var userTeamID uint
	h.DB.Raw("SELECT team_id FROM users WHERE id = ?", claims.ID).Scan(&userTeamID)
	if userTeamID != 0 {
		c.AbortWithStatusJSON(403, gin.H{"error": "вы уже являетесь участником команды перевода"})
		return
	}

	desiredTeamID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id команды должен быть числом"})
		return
	}

	var requestBody struct {
		IntroductoryMessage string `json:"introductoryMessage"`
		RoleID              uint   `json:"roleId"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.AbortWithStatusJSON(404, gin.H{"error": err.Error()})
		return
	}

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var existingTeamID uint
	tx.Raw("SELECT id FROM teams WHERE id = ?", desiredTeamID).Scan(&existingTeamID)
	if existingTeamID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "команда перевода не найдена"})
		return
	}

	if requestBody.RoleID != 0 {
		var existingRoleID uint
		tx.Raw("SELECT id FROM roles WHERE id = ?", requestBody.RoleID).Scan(&existingRoleID)
		if existingRoleID == 0 {
			c.AbortWithStatusJSON(404, gin.H{"error": "роль не найдена"})
			return
		}
	}

	var userRequestToThisTeamID uint
	h.DB.Raw("SELECT id FROM team_join_requests WHERE candidate_id = ? AND team_id = ?", claims.ID, existingTeamID).Scan(&userRequestToThisTeamID)
	if userRequestToThisTeamID != 0 {
		c.AbortWithStatusJSON(409, gin.H{"error": "у вас уже есть заявка на вступление в эту команду"})
		return
	}

	joinRequest := models.TeamJoinRequest{
		CandidateID:         claims.ID,
		TeamID:              existingTeamID,
		IntroductoryMessage: requestBody.IntroductoryMessage,
	}
	if requestBody.RoleID != 0 {
		joinRequest.RoleID = sql.NullInt64{Int64: int64(requestBody.RoleID), Valid: true}
	}

	if result := tx.Create(&joinRequest); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "заявка на вступление в команду успешно отправлена"})
	// Уведомление лидеру команды
}
