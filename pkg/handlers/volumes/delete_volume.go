package volumes

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	"github.com/gin-gonic/gin"
)

func (h handler) DeleteVolume(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	volumeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id тома"})
		return
	}

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var doesVolumeExist bool

	if err := tx.Raw(
		`SELECT EXISTS(
			SELECT 1 FROM titles AS t
			INNER JOIN volumes AS v ON v.title_id = t.id
			WHERE v.id = ? AND t.team_id = (SELECT team_id FROM users WHERE id = ?) 
		)`, volumeID, claims.ID,
	).Scan(&doesVolumeExist).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if !doesVolumeExist {
		c.AbortWithStatusJSON(404, gin.H{"error": "том не найден среди томов тайтлов, переводиммых вашей командой"})
		return
	}

	result := tx.Exec("DELETE FROM volumes WHERE id = ?", volumeID)

	if result.Error != nil {
		if dbErrors.IsForeignKeyViolation(result.Error, constraints.FkChaptersVolume) {
			c.AbortWithStatusJSON(409, gin.H{"error": "удалить можно только том без глав"})
		} else {
			log.Println(result.Error)
			c.AbortWithStatusJSON(409, gin.H{"error": result.Error.Error()})
		}
		return
	}

	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(500, gin.H{"error": "произошла ошибка при удалении тома"})
		return
	}

	tx.Commit()

	c.JSON(200, gin.H{"success": "том успешно удален"})
}
