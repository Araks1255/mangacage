package moderation

import (
	"context"
	"errors"

	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	mongoModels "github.com/Araks1255/mangacage/pkg/common/models/mongo"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type CreateUserOnModerationOptions struct {
	ExistingID     uint
	ProfilePicture []byte
	Collection     *mongo.Collection
}

func CreateUserOnModeration(db *gorm.DB, opts ...CreateUserOnModerationOptions) (uint, error) {
	if len(opts) > 1 {
		return 0, errors.New("объектов опций не может быть больше одного")
	}

	userName := uuid.New().String()
	user := models.UserOnModeration{
		UserName: &userName,
	}

	if len(opts) != 0 {
		if opts[0].ExistingID != 0 {
			user.ExistingID = &opts[0].ExistingID
		}
	}

	tx := db.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	if result := tx.Create(&user); result.Error != nil {
		return 0, result.Error
	}

	if len(opts) == 0 {
		tx.Commit()
		return user.ID, nil
	}

	if opts[0].ProfilePicture != nil {
		if opts[0].Collection == nil {
			return 0, errors.New("передана аватарка, но не передана коллекция")
		}

		userOnModerationProfilePicture := mongoModels.UserOnModerationProfilePicture{
			UserOnModerationID: user.ID,
			ProfilePicture:     opts[0].ProfilePicture,
		}

		if user.ExistingID != nil {
			userOnModerationProfilePicture.CreatorID = *user.ExistingID
		}

		if _, err := opts[0].Collection.InsertOne(context.Background(), userOnModerationProfilePicture); err != nil {
			return 0, err
		}
	}

	tx.Commit()

	return user.ID, nil
}
