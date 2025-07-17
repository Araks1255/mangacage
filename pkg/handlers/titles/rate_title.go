package titles

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	"github.com/gin-gonic/gin"
)

func (h handler) RateTitle(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	titleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id тайтла"})
		return
	}

	var requestBody struct {
		Rate int `json:"rate" binding:"required,gte=1,lte=5"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	result := h.DB.Exec(
		`INSERT INTO title_rates (title_id, user_id, rate)
		SELECT
			t.id, uvc.user_id, ?
		FROM
			titles AS t
			INNER JOIN chapters AS c ON c.title_id = t.id
			INNER JOIN user_viewed_chapters AS uvc ON uvc.chapter_id = c.id AND uvc.user_id = ?
		WHERE
			t.id = ?
		GROUP BY
			t.id, uvc.user_id
		HAVING
			COUNT(DISTINCT c.id) >= 10
		OR
			COUNT(DISTINCT c.id) = (
				SELECT COUNT(chapters.id)
				FROM chapters
				WHERE chapters.title_id = t.id
			)
			
		ON CONFLICT (title_id, user_id)
		DO UPDATE SET rate = EXCLUDED.rate`,
		requestBody.Rate, claims.ID, titleID,
	)

	if result.Error != nil {
		if dbErrors.IsForeignKeyViolation(result.Error, constraints.FkTitleRatesTitle) {
			c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден"})
		} else {
			log.Println(result.Error)
			c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		}
		return
	}

	if result.RowsAffected == 0 {
		c.AbortWithStatusJSON(409, gin.H{"error": "вы прочитали недостаточное количество глав для оценки этого тайтла"})
		return
	}

	c.JSON(201, gin.H{"success": "тайтл успешно оценен"})
}
