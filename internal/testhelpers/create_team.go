package testhelpers

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type CreateTeamOptions struct {
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

	if opts[0].ModeratorID != 0 {
		team.ModeratorID = sql.NullInt64{Int64: int64(opts[0].ModeratorID), Valid: true}
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

	var teamCover struct {
		TeamID uint   `bson:"team_id"`
		Cover  []byte `bson:"cover"`
	}

	teamCover.TeamID = team.ID
	teamCover.Cover = opts[0].Cover

	if _, err := opts[0].Collection.InsertOne(context.Background(), teamCover); err != nil {
		return 0, err
	}

	tx.Commit()

	return team.ID, nil
}
