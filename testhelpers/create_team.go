package testhelpers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateTeamOptions struct {
	Description    string
	Cover          []byte
	PathToMediaDir string
	ModeratorID    uint
}

func CreateTeam(db *gorm.DB, creatorID uint, opts ...CreateTeamOptions) (uint, error) {
	if len(opts) > 1 {
		return 0, errors.New("Объектов опций не может быть больше одного")
	}

	team := models.Team{
		Name:      uuid.New().String(),
		CreatorID: &creatorID,
	}

	if len(opts) != 0 {
		if opts[0].Description != "" {
			team.Description = opts[0].Description
		}
		if opts[0].ModeratorID != 0 {
			team.ModeratorID = &opts[0].ModeratorID
		}
	}

	tx := db.Begin()
	defer tx.Rollback()

	if err := db.Create(&team).Error; err != nil {
		return 0, err
	}

	if len(opts) != 0 && len(opts[0].Cover) != 0 {
		if opts[0].PathToMediaDir == "" {
			return 0, errors.New("не передана директория для сохранения медиафайлов")
		}

		team.CoverPath = fmt.Sprintf("%s/teams/%d.jpg", opts[0].PathToMediaDir, team.ID)

		if err := os.MkdirAll(filepath.Dir(team.CoverPath), 0644); err != nil {
			return 0, err
		}

		if err := os.WriteFile(team.CoverPath, opts[0].Cover, 0644); err != nil {
			return 0, err
		}

		if err := tx.Updates(&team).Error; err != nil {
			return 0, err
		}
	}

	tx.Commit()

	return team.ID, nil
}
