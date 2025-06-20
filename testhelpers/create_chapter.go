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

type CreateChapterOptions struct {
	Pages       [][]byte
	Collection  *mongo.Collection
	ModeratorID uint
}

func CreateChapter(db *gorm.DB, volumeID, teamID, creatorID uint, opts ...CreateChapterOptions) (uint, error) {
	if len(opts) > 1 {
		return 0, errors.New("Объектов опций не может быть больше одного")
	}

	chapter := models.Chapter{
		Name:      uuid.New().String(),
		VolumeID:  volumeID,
		CreatorID: creatorID,
		TeamID:    teamID,
	}

	if len(opts) != 0 {
		if opts[0].ModeratorID != 0 {
			chapter.ModeratorID = &opts[0].ModeratorID
		}
		if len(opts[0].Pages) != 0 {
			chapter.NumberOfPages = len(opts[0].Pages)
		}
	}

	tx := db.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	if result := tx.Create(&chapter); result.Error != nil {
		return 0, result.Error
	}

	if len(opts) == 0 || len(opts[0].Pages) == 0 {
		tx.Commit()
		return chapter.ID, nil
	}

	if opts[0].Collection == nil {
		return 0, errors.New("Переданы страницы, но не передана коллекция")
	}

	chapterPages := mongoModels.ChapterPages{
		ChapterID: chapter.ID,
		CreatorID: creatorID,
		Pages:     opts[0].Pages,
	}

	if _, err := opts[0].Collection.InsertOne(context.Background(), chapterPages); err != nil {
		return 0, err
	}

	tx.Commit()

	return chapter.ID, nil
}
