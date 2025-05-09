package favorites

import (
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	"github.com/gin-gonic/gin"
)

func (h handler) AddChapterToFavorites(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	chapterID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id главы"})
		return
	}

	err = h.DB.Exec("INSERT INTO user_favorite_chapters (user_id, chapter_id) VALUES (?, ?)", claims.ID, chapterID).Error

	if err != nil {
		if dbErrors.IsForeignKeyViolation(err, constraints.FkUserFavoriteChaptersChapter) {
			c.AbortWithStatusJSON(404, gin.H{"error": "глава не найдена"})
			return
		}

		if dbErrors.IsUniqueViolation(err, constraints.UserFavoriteChaptersPkey) {
			c.AbortWithStatusJSON(409, gin.H{"error": "глава уже есть в вашем избранном"})
			return
		}

		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"success": "глава успешно добавлена в ваше избранное"})
}
