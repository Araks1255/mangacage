package seeder

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func seedChapters(ctx context.Context, db *gorm.DB, chaptersPages *mongo.Collection) error {
	var chapterID uint
	if result := db.Raw(
		`INSERT INTO chapters (name, description, number_of_pages, volume_id, creator_id, moderator_id)
		SELECT 'chapter_test', '', 1, v.id, u.id, u.id FROM
		(SELECT id FROM volumes WHERE name = 'volume_test') AS v,
		(SELECT id FROM users WHERE user_name = 'user_test') AS u
		ON CONFLICT DO NOTHING
		RETURNING id`,
	).Scan(&chapterID); result.Error != nil {
		return result.Error
	}

	if chapterID == 0 {
		log.Println("Глава уже создана")
		return nil
	}

	var chapterPages struct {
		ChapterID uint     `bson:"chapter_id"`
		Pages     [][]byte `bson:"pages"`
	}

	chapterPages.ChapterID = chapterID
	chapterPages.Pages = make([][]byte, 1, 1)

	var err error
	chapterPages.Pages[0], err = os.ReadFile("test_data/chapter_page.png")
	if err != nil {
		return err
	}

	if _, err = chaptersPages.InsertOne(ctx, chapterPages); err != nil {
		return err
	}

	return nil
}
