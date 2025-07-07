package moderation

import (
	"context"
	"errors"

	"github.com/Araks1255/mangacage/pkg/common/models"
	mongoModels "github.com/Araks1255/mangacage/pkg/common/models/mongo"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type CreateTeamOnModerationOptions struct {
	ExistingID uint
	Cover      []byte
	Collection *mongo.Collection
}

func CreateTeamOnModeration(db *gorm.DB, userID uint, opts ...CreateTeamOnModerationOptions) (uint, error) {
	if len(opts) > 1 {
		return 0, errors.New("объектов опций не может быть больше одного")
	}

	name := uuid.New().String()
	team := models.TeamOnModeration{
		Name:      &name,
		CreatorID: userID,
	}

	if len(opts) != 0 && opts[0].ExistingID != 0 {
		team.ExistingID = &opts[0].ExistingID
	}

	if err := db.Create(&team).Error; err != nil {
		return 0, err
	}

	if len(opts) == 0 {
		return team.ID, nil
	}

	if opts[0].Cover != nil {
		if opts[0].Collection == nil {
			return 0, errors.New("передана обложка, но не передана коллекция")
		}

		teamCover := mongoModels.TeamOnModerationCover{
			TeamOnModerationID: team.ID,
			CreatorID:          userID,
			Cover:              opts[0].Cover,
		}

		if _, err := opts[0].Collection.InsertOne(context.Background(), teamCover); err != nil {
			return 0, err
		}
	}

	return team.ID, nil
}
