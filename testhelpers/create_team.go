package testhelpers

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

type CreateTeamOptions struct {
	Description string
	Cover       []byte
	Collection  *mongo.Collection
	ModeratorID uint
}

func CreateTeam(db *gorm.DB, creatorID uint, opts ...CreateTeamOptions) (uint, error) {
	if len(opts) > 1 {
		return 0, errors.New("Объектов опций не может быть больше одного")
	}

	team := models.Team{
		Name:      uuid.New().String(),
		CreatorID: creatorID,
	}

	tx := db.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	if len(opts) == 0 {
		if result := tx.Create(&team); result.Error != nil {
			return 0, result.Error
		}
		tx.Commit()
		return team.ID, nil
	}

	if opts[0].Description != "" {
		team.Description = opts[0].Description
	}
	if opts[0].ModeratorID != 0 {
		team.ModeratorID = &opts[0].ModeratorID
	}

	if result := tx.Create(&team); result.Error != nil {
		return 0, result.Error
	}

	if opts[0].Cover == nil && opts[0].Collection == nil {
		tx.Commit()
		return team.ID, nil
	}

	if opts[0].Collection == nil {
		return 0, errors.New("Передана обложка, но не передана коллекция")
	}

	teamCover := mongoModels.TeamCover{
		TeamID:    team.ID,
		CreatorID: creatorID,
		Cover:     opts[0].Cover,
	}

	if _, err := opts[0].Collection.InsertOne(context.Background(), teamCover); err != nil {
		return 0, err
	}

	tx.Commit()

	return team.ID, nil
}
