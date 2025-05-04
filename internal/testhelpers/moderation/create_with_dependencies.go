package moderation

import (
	"errors"

	"github.com/Araks1255/mangacage/internal/testhelpers"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type CreateChapterOnModerationWithDependenciesOptions struct {
	Edited     bool
	Pages      [][]byte
	Collection *mongo.Collection
}

func CreateChapterOnModerationWithDependencies(db *gorm.DB, userID uint, opts ...CreateChapterOnModerationWithDependenciesOptions) (uint, error) {
	if len(opts) > 1 {
		return 0, errors.New("объектов опций не может быть больше одного")
	}

	authorID, err := testhelpers.CreateAuthor(db)
	if err != nil {
		return 0, err
	}

	titleID, err := testhelpers.CreateTitle(db, userID, authorID)
	if err != nil {
		return 0, err
	}

	volumeID, err := testhelpers.CreateVolume(db, titleID, userID)
	if err != nil {
		return 0, err
	}

	if len(opts) != 0 {
		var chapterID uint

		if opts[0].Edited {
			chapterID, err = testhelpers.CreateChapter(db, volumeID, userID)
			if err != nil {
				return 0, err
			}
		}

		if chapterID != 0 && opts[0].Pages != nil {
			chapterOnModerationID, err := CreateChapterOnModeration(
				db, volumeID, userID, CreateChapterOnModerationOptions{ExistingID: chapterID, Pages: opts[0].Pages, Collection: opts[0].Collection},
			)

			if err != nil {
				return 0, err
			}

			return chapterOnModerationID, nil
		}

		if chapterID != 0 {
			chapterOnModerationID, err := CreateChapterOnModeration(db, volumeID, userID, CreateChapterOnModerationOptions{ExistingID: chapterID})
			if err != nil {
				return 0, err
			}
			return chapterOnModerationID, nil
		}

		if opts[0].Pages != nil {
			chapterOnModerationID, err := CreateChapterOnModeration(db, volumeID, userID, CreateChapterOnModerationOptions{Pages: opts[0].Pages, Collection: opts[0].Collection})
			if err != nil {
				return 0, err
			}
			return chapterOnModerationID, nil
		}
	}

	chapterOnModerationID, err := CreateChapterOnModeration(db, volumeID, userID)
	if err != nil {
		return 0, err
	}

	return chapterOnModerationID, nil
}
