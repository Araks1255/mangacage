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

type CreateChapterOptions struct {
	Pages       [][]byte
	Collection  *mongo.Collection
	ModeratorID uint
}

func CreateChapter(db *gorm.DB, volumeID, creatorID uint, opts ...CreateChapterOptions) (uint, error) {
	if len(opts) > 1 {
		return 0, errors.New("Объектов опций не может быть больше одного")
	}

	chapter := models.Chapter{
		Name:      uuid.New().String(),
		VolumeID:  volumeID,
		CreatorID: creatorID,
	}

	if len(opts) != 0 {
		if opts[0].ModeratorID != 0 {
			chapter.ModeratorID = sql.NullInt64{Int64: int64(opts[0].ModeratorID), Valid: true}
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

	var chapterPages struct {
		ChapterID uint     `bson:"chapter_id"`
		Pages     [][]byte `bson:"pages"`
	}

	chapterPages.ChapterID = chapter.ID
	chapterPages.Pages = opts[0].Pages

	if _, err := opts[0].Collection.InsertOne(context.Background(), chapterPages); err != nil {
		return 0, err
	}

	tx.Commit()

	return chapter.ID, nil
}
