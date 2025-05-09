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

type CreateChapterOnModerationOptions struct {
	Pages      [][]byte
	Collection *mongo.Collection
	ExistingID uint
}

func CreateChapterOnModeration(db *gorm.DB, volumeID, userID uint, opts ...CreateChapterOnModerationOptions) (uint, error) {
	if len(opts) > 1 {
		return 0, errors.New("объект опций может быть только один")
	}

	chapter := models.ChapterOnModeration{
		Name:      sql.NullString{String: uuid.New().String(), Valid: true},
		VolumeID:  volumeID,
		CreatorID: userID,
	}

	if len(opts) != 0 {
		if opts[0].ExistingID != 0 {
			chapter.ExistingID = sql.NullInt64{Int64: int64(opts[0].ExistingID), Valid: true}
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

	if len(opts) != 0 && len(opts[0].Pages) != 0 {
		if opts[0].Collection == nil {
			return 0, errors.New("переданы страницы, но не передана коллекция")
		}

		var chapterPages struct {
			ChapterOnModerationID uint     `bson:"chapter_on_moderation_id"`
			Pages                 [][]byte `bson:"pages"`
		}

		chapterPages.ChapterOnModerationID = chapter.ID
		chapterPages.Pages = opts[0].Pages

		if _, err := opts[0].Collection.InsertOne(context.Background(), chapterPages); err != nil {
			return 0, err
		}
	}

	tx.Commit()

	return chapter.ID, nil
}
