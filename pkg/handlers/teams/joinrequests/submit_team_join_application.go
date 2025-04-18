package joinrequests

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

func (h handler) SubmitTeamJoinRequest(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

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

	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	defer tx.Rollback()

	var existingTeamID uint
	tx.Raw("SELECT id FROM teams WHERE id = ?", desiredTeamID).Scan(&existingTeamID)
	if existingTeamID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "команда перевода не найдена"})
		return
	}

	var userRequestToThisTeamID uint
	h.DB.Raw("SELECT id FROM team_join_requests WHERE candidate_id = ? AND team_id = ?", claims.ID, existingTeamID).Scan(&userRequestToThisTeamID)
	if userRequestToThisTeamID != 0 {
		c.AbortWithStatusJSON(409, gin.H{"error": "у вас уже есть заявка на вступление в эту команду"})
		return
	}

	var requestBody struct {
		IntroductoryMessage string `json:"introductoryMessage"`
		Role                string `json:"role"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Println(err) // Тут оба значения опциональны, так что запрос может выполнится даже с пустым телом. Поэтому на валидации json запрос не обрывается, даже при ошибке
	}

	application := models.TeamJoinRequest{
		CandidateID:         claims.ID,
		TeamID:              existingTeamID,
		IntroductoryMessage: requestBody.IntroductoryMessage,
		Role:                requestBody.Role, // Тут пишут что хотят, но будут реальные варианты предложены, если выберут один из них - сразу получат при попадании в команду
	}

	if result := tx.Create(&application); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "заявка на вступление в команду успешно отправлена"})
	// Уведомление лидеру команды
}
