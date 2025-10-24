package moderation

import (
	"errors"

	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateChapterOnModerationOptions struct {
	ExistingID uint
	Volume     uint
	Pages      [][]byte
}

func CreateChapterOnModeration(db *gorm.DB, titleID, teamID, userID uint, opts ...CreateChapterOnModerationOptions) (uint, error) {
	if len(opts) > 1 {
		return 0, errors.New("объект опций может быть только один")
	}

	name := uuid.New().String()

	chapter := models.ChapterOnModeration{
		Name:      &name,
		TitleID:   &titleID,
		CreatorID: &userID,
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

	tx.Commit()

	return chapter.ID, nil
}
