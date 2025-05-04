package testhelpers

import (
	"gorm.io/gorm"
)

func CreateTitleWithDependencies(db *gorm.DB, userID uint) (uint, error) {
	authorID, err := CreateAuthor(db)
	if err != nil {
		return 0, err
	}

	titleID, err := CreateTitle(db, userID, authorID)
	if err != nil {
		return 0, err
	}

	return titleID, nil
}

func CreateVolumeWithDependencies(db *gorm.DB, userID uint) (uint, error) {
	teamID, err := CreateTeam(db, userID)
	if err != nil {
		return 0, err
	}

	if err = AddUserToTeam(db, userID, teamID); err != nil {
		return 0, err
	}

	authorID, err := CreateAuthor(db)
	if err != nil {
		return 0, err
	}

	titleID, err := CreateTitle(db, userID, authorID)
	if err != nil {
		return 0, err
	}

	if err = TranslateTitle(db, teamID, titleID); err != nil {
		return 0, err
	}

	volumeID, err := CreateVolume(db, titleID, userID)
	if err != nil {
		return 0, err
	}

	return volumeID, nil
}

func CreateChapterWithDependencies(db *gorm.DB, userID uint) (uint, error) {
	authorID, err := CreateAuthor(db)
	if err != nil {
		return 0, err
	}

	titleID, err := CreateTitle(db, userID, authorID)
	if err != nil {
		return 0, err
	}

	volumeID, err := CreateVolume(db, titleID, userID)
	if err != nil {
		return 0, err
	}

	chapterID, err := CreateChapter(db, volumeID, userID)
	if err != nil {
		return 0, err
	}

	return chapterID, nil
}
