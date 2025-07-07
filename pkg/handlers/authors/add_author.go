package authors

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/handlers/helpers"
	"github.com/Araks1255/mangacage/pkg/handlers/helpers/authors"
	"github.com/gin-gonic/gin"
)

func (h handler) AddAuthor(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var requestBody models.AuthorOnModerationDTO

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	newAuthor := requestBody.ToAuthorOnModeration(claims.ID)

	exists, err := helpers.CheckEntityWithTheSameNameExistence(h.DB, "authors", &newAuthor.Name, nil, &newAuthor.OriginalName)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if exists {
		c.AbortWithStatusJSON(409, gin.H{"error": "автор с таким именем уже существует"})
		return
	}

	err = h.DB.Create(&newAuthor).Error

	if err != nil {
		code, err := authors.ParseAuthorOnModerationInsertError(err)
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, gin.H{"success": "автор успешно отправлен на модерацию"})
	// Уведомление
}
