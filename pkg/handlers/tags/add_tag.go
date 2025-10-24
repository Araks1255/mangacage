package tags

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

func (h handler) AddTag(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var requestBody dto.CreateTagDTO

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	exists, err := helpers.CheckEntityWithTheSameNameExistence(h.DB, "tags", &requestBody.Name, nil, nil)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if exists {
		c.AbortWithStatusJSON(409, gin.H{"error": "тег с таким названием уже существует"})
		return
	}

	tag := requestBody.ToTagOnModeration(claims.ID)

	err = h.DB.Create(&tag).Error

	if err != nil {
		if dbErrors.IsUniqueViolation(err, constraints.UniqTagOnModerationName) {
			c.AbortWithStatusJSON(409, gin.H{"error": "тег с таким названием уже ожидает модерации"})
		} else {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(201, gin.H{"success": "тег успешно отправлен на модерацию"})
	
	if _, err := h.NotificationsClient.NotifyAboutNewModerationRequest(
		c.Request.Context(),
		&pb.ModerationRequest{
			EntityOnModeration: enums.EntityOnModeration_ENTITY_ON_MODERATION_TAG,
			ID: uint64(tag.ID),
		},
	); err != nil {
		log.Println(err)
	}
}
