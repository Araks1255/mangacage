package moderation

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type CreateUserOnModerationOptions struct {
	ExistingID     uint
	Roles          []string
	ProfilePicture []byte
	Collection     *mongo.Collection
}

func CreateUserOnModeration(db *gorm.DB, opts ...CreateUserOnModerationOptions) (uint, error) {
	if len(opts) > 1 {
		return 0, errors.New("объектов опций не может быть больше одного")
	}

	user := models.UserOnModeration{
		UserName: sql.NullString{String: uuid.New().String(), Valid: true},
	}

	if len(opts) != 0 {
		if opts[0].ExistingID != 0 {
			user.ExistingID = sql.NullInt64{Int64: int64(opts[0].ExistingID), Valid: true}
		}
		if len(opts[0].Roles) != 0 {
			user.Roles = opts[0].Roles
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

		var userProfilePicture struct {
			UserOnModerationID uint   `bson:"user_on_moderation_id"`
			ProfilePicture     []byte `bson:"profile_picture"`
		}

		userProfilePicture.UserOnModerationID = user.ID
		userProfilePicture.ProfilePicture = opts[0].ProfilePicture

		if _, err := opts[0].Collection.InsertOne(context.Background(), userProfilePicture); err != nil {
			return 0, err
		}
	}

	tx.Commit()

	return user.ID, nil
}
