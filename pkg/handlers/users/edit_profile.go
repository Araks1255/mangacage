package users

import (
	"context"
	"errors"
	"log"
	"mime/multipart"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/Araks1255/mangacage/pkg/common/utils"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	"github.com/Araks1255/mangacage/pkg/handlers/helpers"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

func (h handler) EditProfile(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var requestBody dto.EditUserDTO

	if err := c.ShouldBindWith(&requestBody, binding.FormMultipart); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	code, err := checkEditProfileConflicts(h.DB, requestBody, claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	editedProfile := requestBody.ToUserOnModeration(claims.ID)

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
		err = upsertProfilePicture(c.Request.Context(), h.UsersProfilePictures, requestBody.ProfilePicture, editedProfile.ID, claims.ID)
		if err != nil {
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

func checkEditProfileConflicts(db *gorm.DB, requestBody dto.EditUserDTO, userID uint) (code int, err error) {
	ok, err := utils.HasAnyNonEmptyFields(&requestBody)
	if err != nil {
		return 500, err
	}

	if !ok {
		return 400, errors.New("необходим как минимум 1 изменяемый параметр")
	}

	if requestBody.UserName != nil {
		exists, err := helpers.CheckEntityWithTheSameNameExistence(db, "users", requestBody.UserName, nil, nil)
		if err != nil {
			return 500, err
		}

		if exists {
			return 409, errors.New("пользователь с таким именем уже существует")
		}
	}

	return 0, nil
}

func upsertProfilePicture(ctx context.Context, collection *mongo.Collection, pictureFileHeader *multipart.FileHeader, userOnModerationID, userID uint) error {
	profilePicture, err := utils.ReadMultipartFile(pictureFileHeader, 2<<20)
	if err != nil {
		return err
	}

	filter := bson.M{"user_on_moderation_id": userOnModerationID}
	update := bson.M{"$set": bson.M{"profilePicture": profilePicture, "creator_id": userID}}
	opts := options.Update().SetUpsert(true)

	if _, err := collection.UpdateOne(ctx, filter, update, opts); err != nil {
		return err
	}

	return nil
}
