package genres

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	"github.com/Araks1255/mangacage/pkg/handlers/helpers"
	"github.com/gin-gonic/gin"
	"github.com/Araks1255/mangacage_protos/gen/enums"
	pb "github.com/Araks1255/mangacage_protos/gen/site_notifications"
)

func (h handler) AddGenre(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var requestBody dto.CreateGenreDTO

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	exists, err := helpers.CheckEntityWithTheSameNameExistence(h.DB, "genres", &requestBody.Name, nil, nil)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if exists {
		c.AbortWithStatusJSON(409, gin.H{"error": "жанр с таким названием уже существует"})
		return
	}

	genre := requestBody.ToGenreOnModeration(claims.ID)

	err = h.DB.Create(&genre).Error

	if err != nil {
		if dbErrors.IsUniqueViolation(err, constraints.UniqGenreOnModerationName) {
			c.AbortWithStatusJSON(409, gin.H{"error": "жанр с таким названием уже ожидает модерации"})
		} else {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(201, gin.H{"success": "жанр успешно отправлен на модерацию"})
	
	if _, err := h.NotificationsClient.NotifyAboutNewModerationRequest(
		c.Request.Context(),
		&pb.ModerationRequest{
			EntityOnModeration: enums.EntityOnModeration_ENTITY_ON_MODERATION_GENRE,
			ID: uint64(genre.ID),
		},
	); err != nil {
		log.Println(err)
	}
}
