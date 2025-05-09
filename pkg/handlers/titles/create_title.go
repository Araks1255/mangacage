package titles

import (
	"database/sql"
	"log"
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
)

func (h handler) CreateTitle(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	if len(form.Value["name"]) == 0 || len(form.Value["authorId"]) == 0 || len(form.Value["genresIds"]) == 0 && len(form.File["cover"]) == 0 {
		c.AbortWithStatusJSON(400, gin.H{"error": "в запросе недостаточно данных"})
		return
	}

	name := form.Value["name"][0]
	genresIDs := form.Value["genresIds"]

	coverFileHeader := form.File["cover"][0]
	if coverFileHeader.Size > 10<<20 {
		c.AbortWithStatusJSON(400, gin.H{"error": "слишком большой размер фото (лимит 10мб)"})
		return
	}

	var description string
	if len(form.Value["description"]) != 0 {
		description = form.Value["description"][0]
	}

	desiredAuthorID, err := strconv.ParseUint(form.Value["authorId"][0], 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id автора"})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var doesTitleWithTheSameNameExist bool
	if err := tx.Raw("SELECT EXISTS(SELECT 1 FROM titles WHERE lower(name) = lower(?))", name).Scan(&doesTitleWithTheSameNameExist).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if doesTitleWithTheSameNameExist {
		c.AbortWithStatusJSON(409, gin.H{"error": "тайтл с таким названием уже существует"})
		return
	}

	newTitle := models.TitleOnModeration{
		Name:        sql.NullString{String: name, Valid: true},
		Description: description,
		CreatorID:   claims.ID,
		AuthorID:    sql.NullInt64{Int64: int64(desiredAuthorID), Valid: true},
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
		SELECT ?, UNNEST(?::INTEGER[])`,
		newTitle.ID, pq.Array(genresIDs),
	).Error

	if err != nil {
		if dbErrors.IsForeignKeyViolation(err, constraints.FkTitleGenresGenre) {
			c.AbortWithStatusJSON(409, gin.H{"error": "указаны невалидные жанры"})
		} else {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}
		return
	}

	var titleCover struct {
		TitleOnModerationID uint   `bson:"title_on_moderation"`
		Cover               []byte `bson:"cover"`
	}

	titleCover.TitleOnModerationID = newTitle.ID
	titleCover.Cover, err = utils.ReadMultipartFile(coverFileHeader, 10<<20)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if _, err = h.TitlesOnModerationCovers.InsertOne(c.Request.Context(), titleCover); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"error": "тайтл успешно отправлен на модерацию"})

	if _, err = h.NotificationsClient.NotifyAboutTitleOnModeration(c.Request.Context(), &pb.TitleOnModeration{ID: uint64(newTitle.ID), New: true}); err != nil {
		log.Println(err)
	}
}
