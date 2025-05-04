package testhelpers

import (
	"errors"

	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type CreateAuthorOptions struct {
	Genres []string
}

func CreateAuthor(db *gorm.DB, opts ...CreateAuthorOptions) (uint, error) {
	if len(opts) > 1 {
		return 0, errors.New("Объектов опций не может быть больше одного")
	}

	author := models.Author{
		Name: uuid.New().String(),
	}

	tx := db.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	if result := tx.Create(&author); result.Error != nil {
		return 0, result.Error
	}

	if len(opts) == 0 {
		tx.Commit()
		return author.ID, nil
	}

	if result := tx.Exec(
		`INSERT INTO author_genres (author_id, genre_id)
		SELECT ?, genres.id FROM genres
		JOIN UNNEST(?::TEXT[]) AS genre_name ON genres.name = genre_name`,
		author.ID, pq.Array(opts[0].Genres),
	); result.Error != nil {
		return 0, result.Error
	}

	tx.Commit()

	return author.ID, nil
}
