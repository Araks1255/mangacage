package moderation

import (
	"errors"

	"github.com/Araks1255/mangacage/testhelpers"
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

	titleID, err := testhelpers.CreateTitleWithDependencies(db, userID)
	if err != nil {
		return 0, err
	}

	teamID, err := testhelpers.CreateTeam(db, userID)
	if err != nil {
		return 0, err
	}

	volumeID, err := testhelpers.CreateVolume(db, titleID, teamID, userID)
	if err != nil {
		return 0, err
	}

	if len(opts) != 0 {
		var chapterID uint

		if opts[0].Edited {
			chapterID, err = testhelpers.CreateChapter(db, volumeID, teamID, userID)
			if err != nil {
				return 0, err
			}
		}

		if chapterID != 0 && opts[0].Pages != nil {
			chapterOnModerationID, err := CreateChapterOnModeration(
				db, volumeID, teamID, userID, CreateChapterOnModerationOptions{ExistingID: chapterID, Pages: opts[0].Pages, Collection: opts[0].Collection},
			)

			if err != nil {
				return 0, err
			}

			return chapterOnModerationID, nil
		}

		if chapterID != 0 {
			chapterOnModerationID, err := CreateChapterOnModeration(db, volumeID, teamID, userID, CreateChapterOnModerationOptions{ExistingID: chapterID})
			if err != nil {
				return 0, err
			}
			return chapterOnModerationID, nil
		}

		if opts[0].Pages != nil {
			chapterOnModerationID, err := CreateChapterOnModeration(db, volumeID, teamID, userID, CreateChapterOnModerationOptions{Pages: opts[0].Pages, Collection: opts[0].Collection})
			if err != nil {
				return 0, err
			}
			return chapterOnModerationID, nil
		}
	}

	chapterOnModerationID, err := CreateChapterOnModeration(db, volumeID, teamID, userID)
	if err != nil {
		return 0, err
	}

	return chapterOnModerationID, nil
}

func CreateVolumeOnModerationWithDependencies(db *gorm.DB, userID uint, edited bool) (uint, error) {
	titleID, err := testhelpers.CreateTitleWithDependencies(db, userID)
	if err != nil {
		return 0, err
	}

	teamID, err := testhelpers.CreateTeam(db, userID)
	if err != nil {
		return 0, err
	}

	var volumeID uint

	if edited {
		existingVolumeID, err := testhelpers.CreateVolumeWithDependencies(db, userID)
		if err != nil {
			return 0, err
		}

		volumeID, err = CreateVolumeOnModeration(db, titleID, teamID, userID, CreateVolumeOnModerationOptions{ExistingID: existingVolumeID})
	} else {
		volumeID, err = CreateVolumeOnModeration(db, titleID, teamID, userID)
	}

	if err != nil {
		return 0, err
	}

	return volumeID, nil
}
