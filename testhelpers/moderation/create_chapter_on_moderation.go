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

type CreateChapterOnModerationOptions struct {
	ExistingID uint
	Volume     uint
	Pages      [][]byte
	Collection *mongo.Collection
}

func CreateChapterOnModeration(db *gorm.DB, titleID, teamID, userID uint, opts ...CreateChapterOnModerationOptions) (uint, error) {
	if len(opts) > 1 {
		return 0, errors.New("объект опций может быть только один")
	}

	name := uuid.New().String()
	chapter := models.ChapterOnModeration{
		Name:      &name,
		TitleID:   &titleID,
		CreatorID: userID,
		TeamID:    teamID,
	}

	if len(opts) != 0 {
		if opts[0].ExistingID != 0 {
			chapter.ExistingID = &opts[0].ExistingID
		}
		if opts[0].Volume != 0 {
			chapter.Volume = &opts[0].Volume
		}
		if len(opts[0].Pages) != 0 {
			numberOfPages := len(opts[0].Pages)
			chapter.NumberOfPages = &numberOfPages
		}
	}

	tx := db.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	if result := tx.Create(&chapter); result.Error != nil {
		return 0, result.Error
	}

	if len(opts) != 0 && len(opts[0].Pages) != 0 {
		if opts[0].Collection == nil {
			return 0, errors.New("переданы страницы, но не передана коллекция")
		}

		chapterPages := mongoModels.ChapterOnModerationPages{
			ChapterOnModerationID: chapter.ID,
			CreatorID:             userID,
			Pages:                 opts[0].Pages,
		}

		if _, err := opts[0].Collection.InsertOne(context.Background(), chapterPages); err != nil {
			return 0, err
		}
	}

	tx.Commit()

	return chapter.ID, nil
}
