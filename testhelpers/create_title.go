package testhelpers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type CreateTitleOptions struct {
	Description    string
	Cover          []byte
	PathToMediaDir string
	TeamID         uint
	Views          uint
	ModeratorID    uint
	Genres         []string
	Tags           []string
}

// только jpg
func CreateTitle(db *gorm.DB, creatorID, authorID uint, opts ...CreateTitleOptions) (uint, error) {
	if len(opts) > 1 {
		return 0, errors.New("Объектов опций не может быть больше одного")
	}

	title := models.Title{
		Name:              uuid.New().String(),
		EnglishName:       uuid.New().String(),
		OriginalName:      uuid.New().String(),
		PublishingStatus:  "ongoing",
		TranslatingStatus: "ongoing",
		AgeLimit:          18,
		YearOfRelease:     1999,
		Type:              "manga",
		AuthorID:          authorID,
		CreatorID:         &creatorID,
	}

	tx := db.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	if len(opts) == 0 {
		if result := tx.Create(&title); result.Error != nil {
			return 0, result.Error
		}
		tx.Commit()
		return title.ID, nil
	}

	if opts[0].Description != "" {
		title.Description = &opts[0].Description
	}
	if opts[0].ModeratorID != 0 {
		title.ModeratorID = &opts[0].ModeratorID
	}
	if opts[0].Views != 0 {
		title.Views = opts[0].Views
	}

	if result := tx.Create(&title); result.Error != nil {
		return 0, result.Error
	}

	if opts[0].TeamID != 0 {
		err := tx.Exec("INSERT INTO title_teams(title_id, team_id) VALUES (?, ?)", title.ID, opts[0].TeamID).Error
		if err != nil {
			return 0, err
		}
	}

	if opts[0].Genres != nil {
		if result := tx.Exec(
			`INSERT INTO title_genres (title_id, genre_id)
			SELECT ?, genres.id FROM genres
			JOIN UNNEST(?::TEXT[]) AS genre_name ON genres.name = genre_name`,
			title.ID, pq.Array(opts[0].Genres),
		); result.Error != nil {
			return 0, result.Error
		}
	}

	if opts[0].Tags != nil {
		if result := tx.Exec(
			`INSERT INTO title_tags (title_id, tag_id)
			SELECT ?, tags.id FROM tags
			JOIN UNNEST(?::TEXT[]) AS tag_name ON tags.name = tag_name`,
			title.ID, pq.Array(opts[0].Tags),
		); result.Error != nil {
			return 0, result.Error
		}
	}

	if opts[0].Cover != nil {
		if opts[0].PathToMediaDir == "" {
			return 0, errors.New("не передана директория для сохранения медиафайлов")
		}

		title.CoverPath = fmt.Sprintf("%s/titles/%d.jpg", opts[0].PathToMediaDir, title.ID)

		if err := os.MkdirAll(filepath.Dir(title.CoverPath), 0644); err != nil {
			return 0, err
		}

		if err := os.WriteFile(title.CoverPath, opts[0].Cover, 0644); err != nil {
			return 0, err
		}

		if err := tx.Updates(&title).Error; err != nil {
			return 0, err
		}
	}

	tx.Commit()

	return title.ID, nil
}
