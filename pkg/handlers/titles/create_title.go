package titles

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"mime/multipart"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/common/utils"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

func (h handler) CreateTitle(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	form, err := c.MultipartForm()
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	name, description, authorID, genresIDs, coverFileHeader, err := parseCreateTitleParams(form)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	ok, err := checkTitleWithTheSameNameExistence(tx, name)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if ok {
		c.AbortWithStatusJSON(409, gin.H{"error": "тайтл с таким названием уже существует"})
		return
	}

	newTitle := models.TitleOnModeration{
		Name:        sql.NullString{String: name, Valid: true},
		Description: description,
		CreatorID:   claims.ID,
		AuthorID:    &authorID,
	}

	err = tx.Create(&newTitle).Error

	if err != nil {
		if dbErrors.IsUniqueViolation(err, constraints.UniTitlesOnModerationName) {
			c.AbortWithStatusJSON(409, gin.H{"error": "тайтл с таким названием уже ожидает модерации"})
			return
		}

		if dbErrors.IsForeignKeyViolation(err, constraints.FkTitlesOnModerationAuthor) {
			c.AbortWithStatusJSON(404, gin.H{"error": "автор не найден"})
			return
		}

		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	err = tx.Exec(
		`INSERT INTO title_on_moderation_genres (title_on_moderation_id, genre_id)
		SELECT ?, UNNEST(?::BIGINT[])`,
		newTitle.ID, pq.Array(genresIDs),
	).Error

	if err != nil {
		if dbErrors.IsForeignKeyViolation(err, constraints.FkTitleOnModerationGenresGenre) {
			c.AbortWithStatusJSON(404, gin.H{"error": "жанры не найдены"})
		} else {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}
		return
	}

	if err := insertTitleCover(c.Request.Context(), h.TitlesOnModerationCovers, newTitle.ID, coverFileHeader); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "тайтл успешно отправлен на модерацию"})

	if _, err = h.NotificationsClient.NotifyAboutTitleOnModeration(c.Request.Context(), &pb.TitleOnModeration{ID: uint64(newTitle.ID), New: true}); err != nil {
		log.Println(err)
	}
}

func parseCreateTitleParams(form *multipart.Form) (name, description string, authorID uint, genresIDs []uint, coverFileHeader *multipart.FileHeader, err error) {
	if len(form.Value["name"]) == 0 || len(form.Value["authorId"]) == 0 || len(form.Value["genresIds"]) == 0 || len(form.File["cover"]) == 0 {
		return "", "", 0, nil, nil, errors.New("в запросе недостаточно данных")
	}

	genresIDs, err = utils.ParseUintSlice(form.Value["genresIds"])
	if err != nil {
		return "", "", 0, nil, nil, errors.New("указаны невалидные id жанров")
	}

	authorIDuint64, err := strconv.ParseUint(form.Value["authorId"][0], 10, 64)
	if err != nil {
		return "", "", 0, nil, nil, errors.New("указан невалидный id автора")
	}

	coverFileHeader = form.File["cover"][0]
	if coverFileHeader.Size > 2<<20 {
		return "", "", 0, nil, nil, errors.New("превышен лимит размера обложки (2мб)")
	}

	if len(form.Value["description"]) != 0 {
		description = form.Value["description"][0]
	}

	return form.Value["name"][0], description, uint(authorIDuint64), genresIDs, coverFileHeader, nil
}

func checkTitleWithTheSameNameExistence(db *gorm.DB, name string) (bool, error) {
	var doesTitleWithTheSameNameExist bool

	if err := db.Raw("SELECT EXISTS(SELECT 1 FROM titles WHERE lower(name) = lower(?))", name).Scan(&doesTitleWithTheSameNameExist).Error; err != nil {
		return false, err
	}

	if doesTitleWithTheSameNameExist {
		return true, nil
	}

	return false, nil
}

func insertTitleCover(ctx context.Context, collection *mongo.Collection, titleID uint, coverFileHeader *multipart.FileHeader) (err error) {
	var titleCover struct {
		TitleOnModerationID uint   `bson:"title_on_moderation_id"`
		Cover               []byte `bson:"cover"`
	}

	titleCover.TitleOnModerationID = titleID
	titleCover.Cover, err = utils.ReadMultipartFile(coverFileHeader, 2<<20)

	if err != nil {
		return err
	}

	if _, err = collection.InsertOne(ctx, titleCover); err != nil {
		return err
	}

	return nil
}
