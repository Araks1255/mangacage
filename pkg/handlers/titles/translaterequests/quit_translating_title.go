package translaterequests

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"

	"github.com/gin-gonic/gin"
)

func (h handler) QuitTranslatingTitle(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	titleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id тайтла должен быть числом"})
		return
	}

	result := h.DB.Exec("DELETE FROM title_teams WHERE title_id = ? AND team_id = (SELECT team_id FROM users WHERE id = ?)", titleID, claims.ID)

	if result.Error != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден среди переводимых вашей командой"})
		return
	}

	c.JSON(200, gin.H{"success": "ваша команда больше не переводит этот тайтл"})
}
