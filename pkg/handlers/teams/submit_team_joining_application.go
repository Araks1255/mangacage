package teams

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) SubmitTeamJoiningRequest(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var userTeamID uint
	h.DB.Raw("SELECT team_id FROM users WHERE id = ?", claims.ID).Scan(&userTeamID)
	if userTeamID != 0 {
		c.AbortWithStatusJSON(403, gin.H{"error": "вы уже являетесь участником команды"})
		return
	}

	team := c.Param("team")

	var teamID uint
	h.DB.Raw("SELECT id FROM teams WHERE name = ?", team).Scan(&teamID)
	if teamID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "команда не найдена"})
		return
	}

	var userApplicationToThisTeamID uint
	h.DB.Raw("SELECT id FROM team_joining_applications WHERE candidate_id = ? AND team_id = ?", claims.ID, teamID).Scan(&userApplicationToThisTeamID)
	if userApplicationToThisTeamID != 0 {
		c.AbortWithStatusJSON(409, gin.H{"error": "у вас уже есть заявка на вступление в эту команду"})
		return
	}

	var requestBody struct {
		IntroductoryMessage string `json:"introductoryMessage"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Println(err)
	}

	application := models.TeamJoiningApplication{
		CandidateID:         claims.ID,
		TeamID:              teamID,
		IntroductoryMessage: requestBody.IntroductoryMessage,
	}

	if result := h.DB.Create(&application); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	c.JSON(200, gin.H{"success": "заявка на вступление в команду успешно отправлена"})
	// Уведомление лидеру команды
}
