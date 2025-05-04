package titles

import (
	"context"
	"database/sql"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/common/utils"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createTitleBody struct {
	Name        string
	Description string
	AuthorID    uint
	Genres      []string
	Cover       []byte
}

func (h handler) CreateTitle(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	if len(form.Value["name"]) == 0 || len(form.Value["authorId"]) == 0 || len(form.Value["genres"]) == 0 && len(form.File["cover"]) == 0 {
		c.AbortWithStatusJSON(400, gin.H{"error": "в запросе недостаточно данных"})
		return
	}

	name := form.Value["name"][0]
	genres := form.Value["genres"]

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

	var numberOfExistingGenres int64
	h.DB.Raw(
		`SELECT COUNT(*) FROM GENRES
		WHERE name IN
			(SELECT unnest(?::TEXT[]))`,
		pq.Array(genres),
	).Scan(&numberOfExistingGenres)

	if numberOfExistingGenres != int64(len(genres)) {
		c.AbortWithStatusJSON(409, gin.H{"error": "указаны невалидные жанры"})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var existing struct {
		AuthorID            uint
		TitleOnModerationID uint
		TitleID             uint
	}

	tx.Raw(
		`SELECT
			(SELECT id FROM authors WHERE id = ?) AS author_id,
			(SELECT id FROM titles_on_moderation WHERE lower(name) = lower(?)) AS title_on_moderation_id,
			(SELECT id FROM titles WHERE lower(name) = lower(?)) AS title_id`,
		desiredAuthorID, name, name,
	).Scan(&existing)

	if existing.AuthorID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "автор не найден"})
		return
	}
	if existing.TitleOnModerationID != 0 {
		c.AbortWithStatusJSON(409, gin.H{"error": "тайтл с таким названием уже ожидает модерации"})
		return
	}
	if existing.TitleID != 0 {
		c.AbortWithStatusJSON(409, gin.H{"error": "тайтл с таким названием уже существует"})
		return
	}

	newTitle := models.TitleOnModeration{
		Name:        sql.NullString{String: name, Valid: true},
		Description: description,
		CreatorID:   claims.ID,
		AuthorID:    sql.NullInt64{Int64: int64(existing.AuthorID), Valid: true},
		Genres:      genres,
	}

	if result := tx.Create(&newTitle); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	var titleCover struct {
		TitleOnModerationID uint   `bson:"title_on_moderation_id"`
		Cover               []byte `bson:"cover"`
	}

	titleCover.TitleOnModerationID = newTitle.ID
	titleCover.Cover, err = utils.ReadMultipartFile(coverFileHeader, 10<<20)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
	}

	if _, err := h.TitlesOnModerationCovers.InsertOne(context.Background(), titleCover); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "тайтл успешно отправлен на модерацию"})

	if _, err = h.NotificationsClient.NotifyAboutTitleOnModeration(context.Background(), &pb.TitleOnModeration{ID: uint64(newTitle.ID), New: true}); err != nil {
		log.Println(err)
	}
}
