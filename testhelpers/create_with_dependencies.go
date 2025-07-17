package testhelpers

import (
	"gorm.io/gorm"
)

func CreateTitleWithDependencies(db *gorm.DB, userID uint, genres ...string) (uint, error) {
	authorID, err := CreateAuthor(db)
	if err != nil {
		return 0, err
	}

	titleID, err := CreateTitle(db, userID, authorID, CreateTitleOptions{Genres: genres})
	if err != nil {
		return 0, err
	}

	return titleID, nil
}

func CreateChapterWithDependencies(db *gorm.DB, userID uint) (uint, error) {
	titleID, err := CreateTitleWithDependencies(db, userID)
	if err != nil {
		return 0, err
	}
	teamID, err := CreateTeam(db, userID)
	if err != nil {
		return 0, err
	}

	chapterID, err := CreateChapter(db, titleID, teamID, userID)
	if err != nil {
		return 0, err
	}

	return chapterID, nil
}

func CreateTitleTranslatingByUserTeam(db *gorm.DB, userID uint, genres, tags []string) (uint, error) {
	authorID, err := CreateAuthor(db)
	if err != nil {
		return 0, err
	}

	teamID, err := CreateTeam(db, userID)
	if err != nil {
		return 0, err
	}

	if err := AddUserToTeam(db, userID, teamID); err != nil {
		return 0, err
	}

	titleID, err := CreateTitle(db, userID, authorID, CreateTitleOptions{Genres: genres, Tags: tags, TeamID: teamID})
	if err != nil {
		return 0, err
	}

	return titleID, nil
}

func CreateChapterTranslatingByUserTeam(db *gorm.DB, userID uint) (uint, error) {
	teamID, err := CreateTeam(db, userID)
	if err != nil {
		return 0, err
	}

	if err = AddUserToTeam(db, userID, teamID); err != nil {
		return 0, err
	}

	titleID, err := CreateTitleWithDependencies(db, userID)
	if err != nil {
		return 0, err
	}

	if err = TranslateTitle(db, teamID, titleID); err != nil {
		return 0, err
	}

	chapterID, err := CreateChapter(db, titleID, teamID, userID)
	if err != nil {
		return 0, err
	}

	return chapterID, nil
}
