package users

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/common/utils"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	"github.com/Araks1255/mangacage/pkg/handlers/helpers"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h handler) EditProfile(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var requestBody models.UserOnModerationDTO

	c.ShouldBindWith(&requestBody, binding.FormMultipart)

	ok, err := utils.HasAnyNonEmptyFields(&requestBody)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if !ok {
		c.AbortWithStatusJSON(400, gin.H{"error": "необходим как минимум 1 изменяемый параметр"})
		return
	}

	if requestBody.UserName != nil {
		exists, err := helpers.CheckEntityWithTheSameNameExistence(h.DB, "users", requestBody.UserName, nil, nil)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}

		if exists {
			c.AbortWithStatusJSON(409, gin.H{"error": "пользователь с таким именем уже существует"})
			return
		}
	}

	editedProfile := requestBody.ToUserOnModeration(&claims.ID)

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	err = tx.Clauses(helpers.OnExistingIDConflictClause).Create(&editedProfile).Error

	if err != nil {
		if dbErrors.IsUniqueViolation(err, constraints.UniqUserOnModerationUserName) {
			c.AbortWithStatusJSON(409, gin.H{"error": "пользователь с таким именем уже ожидает модерации"})
		} else {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}
		return
	}

	if requestBody.ProfilePicture != nil {
		profilePicture, err := utils.ReadMultipartFile(requestBody.ProfilePicture, 2<<20)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}

		filter := bson.M{"user_on_moderation_id": editedProfile.ID}
		update := bson.M{"$set": bson.M{"profilePicture": profilePicture, "creator_id": claims.ID}}
		opts := options.Update().SetUpsert(true)

		if _, err := h.UsersProfilePictures.UpdateOne(c.Request.Context(), filter, update, opts); err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "изменения профиля успешно отправлены на модерацию"})

	if _, err := h.NotificationsClient.NotifyAboutUserOnModeration(c.Request.Context(), &pb.User{
		ID:  uint64(*editedProfile.ExistingID),
		New: false,
	}); err != nil {
		log.Println(err)
	}
}
