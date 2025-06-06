package favorites

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	"github.com/gin-gonic/gin"
)

func (h handler) AddTitleToFavorites(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	titleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id тайтла"})
		return
	}

	err = h.DB.Exec("INSERT INTO user_favorite_titles (user_id, title_id) VALUES (?, ?)", claims.ID, titleID).Error

	if err != nil {
		if dbErrors.IsForeignKeyViolation(err, constraints.FkUserFavoriteTitlesTitle) {
			c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден"})
			return
		}

		if dbErrors.IsUniqueViolation(err, constraints.UserFavoriteTitlesPkey) {
			c.AbortWithStatusJSON(409, gin.H{"error": "тайтл уже добавлен в ваше избранное"})
			return
		}

		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"success": "тайтл успешно добавлен к вам в избранное"})
}
