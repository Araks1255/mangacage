package moderation

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type CreateTitleOnModerationOptions struct {
	ExistingID uint
	AuthorID   uint
	Genres     []string
	Cover      []byte
	Collection *mongo.Collection
}

func CreateTitleOnModeration(db *gorm.DB, userID uint, opts ...CreateTitleOnModerationOptions) (uint, error) {
	if len(opts) > 1 {
		return 0, errors.New("объектов опций не может быть больше одного")
	}

	title := models.TitleOnModeration{
		Name:      sql.NullString{String: uuid.New().String(), Valid: true},
		CreatorID: userID,
	}

	if len(opts) != 0 {
		if opts[0].ExistingID != 0 {
			title.ExistingID = &opts[0].ExistingID
		}
		if opts[0].AuthorID != 0 {
			title.AuthorID = &opts[0].AuthorID
		}
	}

	tx := db.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	if result := tx.Create(&title); result.Error != nil {
		return 0, result.Error
	}

	if len(opts) == 0 {
		tx.Commit()
		return title.ID, nil
	}

	if opts[0].Genres != nil {
		if err := tx.Exec(
			`INSERT INTO title_on_moderation_genres (title_on_moderation_id, genre_id)
			SELECT ?, genres.id
			FROM genres
			JOIN UNNEST(?::TEXT[]) AS genre_name ON genres.name = genre_name`,
			title.ID, pq.Array(opts[0].Genres),
		).Error; err != nil {
			return 0, err
		}

	}

	if opts[0].Cover != nil {
		if opts[0].Collection == nil {
			return 0, errors.New("передана обложка, но не передана коллекция")
		}

		var titleCover struct {
			TitleOnModerationID uint   `bson:"title_on_moderation_id"`
			Cover               []byte `bson:"cover"`
		}

		titleCover.TitleOnModerationID = title.ID
		titleCover.Cover = opts[0].Cover

		if _, err := opts[0].Collection.InsertOne(context.Background(), titleCover); err != nil {
			return 0, err
		}
	}

	tx.Commit()

	return title.ID, nil
}
