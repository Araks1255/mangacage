package titles

// import (
// 	"database/sql"
// 	"log"

// 	"github.com/Araks1255/mangacage/pkg/common/models"

// 	"github.com/gin-gonic/gin"
// )

// func (h handler) TranslateTitle(c *gin.Context) {
// 	claims := c.MustGet("claims").(*models.Claims)

// 	if claims.Role != "team_owner" {
// 		c.AbortWithStatusJSON(403, gin.H{"error": "Взять тайтл на перевод может только владелец команды"})
// 		return
// 	}

// 	var requestBody struct {
// 		Title string `json:"title" binding:"required"`
// 	}

// 	if err := c.ShouldBindJSON(&requestBody); err != nil {
// 		log.Println(err)
// 		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
// 		return
// 	}

// 	var desiredTitle models.Title
// 	h.DB.Raw("SELECT * FROM titles WHERE name = ?", requestBody.Title).Scan(&desiredTitle)
// 	if desiredTitle.TeamID.Valid {
// 		c.AbortWithStatusJSON(403, gin.H{"error": "Тайтл уже переводит другая команда"})
// 		return
// 	}

// 	var userTeamID sql.NullInt32 //uint
// 	h.DB.Raw("SELECT teams.id FROM teams INNER JOIN users ON teams.id = users.team_id WHERE users.id = ?", claims.ID).Scan(&userTeamID)
// 	if !userTeamID.Valid {
// 		c.AbortWithStatusJSON(403, gin.H{"error": "Вы не состоите в команде перевода"})
// 		return
// 	}

// 	desiredTitle.TeamID = userTeamID

// 	if result := h.DB.Save(&desiredTitle); result.Error != nil {
// 		log.Println(result.Error)
// 		c.AbortWithStatusJSON(500, gin.H{"error": "Не удалось взять тайтл на перевод"})
// 		return
// 	}

// 	c.JSON(201, gin.H{"success": "Теперь ваша команда переводит этот тайтл"})
// }
