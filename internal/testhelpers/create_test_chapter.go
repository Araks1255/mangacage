package testhelpers

import (
	"context"
	"os"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func CreateTestChapter(db *gorm.DB, chaptersPages *mongo.Collection) (id uint, err error) {
	chapter := models.Chapter{
		Name:          uuid.New().String(),
		Description:   "someDescription",
		NumberOfPages: 1,
	}

	row := db.Raw("SELECT volume_id, creator_id, moderator_id FROM chapters WHERE name = 'chapter_test'").Row()

	if err = row.Scan(&chapter.VolumeID, &chapter.CreatorID, &chapter.ModeratorID); err != nil {
		return 0, err
	}

	if result := db.Create(&chapter); result.Error != nil {
		return 0, result.Error
	}

	var chapterPages struct {
		ChapterID uint     `bson:"chapter_id"`
		Pages     [][]byte `bson:"pages"`
	}

	chapterPages.ChapterID = chapter.ID
	chapterPages.Pages = make([][]byte, 1, 1)

	chapterPages.Pages[0], err = os.ReadFile("test_data/chapter_page.png")
	if err != nil {
		return 0, err
	}

	if _, err := chaptersPages.InsertOne(context.Background(), chapterPages); err != nil {
		return 0, err
	}

	return chapter.ID, nil
}
