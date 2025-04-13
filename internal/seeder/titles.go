package seeder

import (
	"context"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func seedTitles(ctx context.Context, db *gorm.DB, titlesCovers *mongo.Collection) error {
	var titleID uint
	if result := db.Raw(
		`INSERT INTO titles (name, description, creator_id, moderator_id, author_id, team_id)
		SELECT 'title_test', '', u.id, u.id, a.id, t.id FROM 
		(SELECT id FROM users WHERE user_name = 'user_test') AS u,
		(SELECT id FROM authors WHERE name = 'author_test') AS a,
		(SELECT id FROM teams WHERE name = 'team_test') AS t
		ON CONFLICT DO NOTHING
		RETURNING id`,
	).Scan(&titleID); result.Error != nil {
		return result.Error
	}

	if titleID == 0 { // Уже создан, ведь returning ничего не вернул
		log.Println("Тайтл уже создан")
		return nil
	}

	var titleCover struct {
		TitleID uint   `bson:"title_id"`
		Cover   []byte `bson:"cover"`
	}

	var err error
	titleCover.TitleID = titleID
	titleCover.Cover, err = os.ReadFile("test_data/title_cover.png")
	if err != nil {
		return err
	}

	if _, err = titlesCovers.InsertOne(ctx, titleCover); err != nil {
		return err
	}

	return nil
}
