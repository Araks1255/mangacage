package testhelpers

import (
	"errors"
	"fmt"
	"os"

	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateChapterOptions struct {
	Pages          [][]byte
	PathToMediaDir string
	Views          uint
	ModeratorID    uint
	Volume         uint
}

// только jpg
func CreateChapter(db *gorm.DB, titleID, teamID, creatorID uint, opts ...CreateChapterOptions) (uint, error) {
	if len(opts) > 1 {
		return 0, errors.New("Объектов опций не может быть больше одного")
	}

	chapter := models.Chapter{
		Name:      uuid.New().String(),
		CreatorID: &creatorID,
		TeamID:    teamID,
		TitleID:   titleID,
	}

	if len(opts) != 0 {
		if opts[0].ModeratorID != 0 {
			chapter.ModeratorID = &opts[0].ModeratorID
		}
		if len(opts[0].Pages) != 0 {
			chapter.NumberOfPages = len(opts[0].Pages)
		}
		if opts[0].Views != 0 {
			chapter.Views = opts[0].Views
		}
		if opts[0].Volume != 0 {
			chapter.Volume = opts[0].Volume
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

	if opts[0].PathToMediaDir == "" {
		return 0, errors.New("не передана директория сохранения медиафайлов")
	}

	pages := make([]models.Page, len(opts[0].Pages))

	chapter.PagesDirPath = fmt.Sprintf("%s/chapters/%d", opts[0].PathToMediaDir, chapter.ID)

	if err := os.MkdirAll(chapter.PagesDirPath, 0755); err != nil {
		return 0, err
	}

	for i := 0; i < len(opts[0].Pages); i++ {
		path := fmt.Sprintf("%s/chapters/%d/%d.jpg", opts[0].PathToMediaDir, chapter.ID, i+1)

		if err := os.WriteFile(path, opts[0].Pages[i], 0755); err != nil {
			return 0, err
		}

		pages[i] = models.Page{
			ChapterID: &chapter.ID,
			Number:    uint(i) + 1,
			Path:      path,
		}
	}

	if err := tx.Create(pages).Error; err != nil {
		return 0, err
	}

	if err := tx.Updates(&chapter).Error; err != nil {
		return 0, err
	}

	tx.Commit()

	return chapter.ID, nil
}
