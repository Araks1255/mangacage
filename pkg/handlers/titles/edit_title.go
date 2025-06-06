package titles

import (
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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (h handler) EditTitle(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	name, description, titleID, authorID, genresIDs, coverFileHeader, err := parseEditTitleParams(form, c.Param)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var doesTitleExist bool

	if err := tx.Raw(
		"SELECT EXISTS(SELECT 1 FROM titles WHERE id = ? AND team_id = (SELECT team_id FROM users WHERE id = ?))",
		titleID, claims.ID,
	).Scan(&doesTitleExist).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if !doesTitleExist {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден среди переводимых вашей командой тайтлов"})
		return
	}

	editedTitle := models.TitleOnModeration{
		ExistingID:  &titleID,
		CreatorID:   claims.ID,
		Description: description,
	}

	if authorID != 0 {
		editedTitle.AuthorID = &authorID
	}

	if name != "" {
		var doesTitleWithTheSameNameExist bool

		if err := tx.Raw("SELECT EXISTS(SELECT 1 FROM titles WHERE lower(name) = lower(?))", form.Value["name"][0]).Scan(&doesTitleWithTheSameNameExist).Error; err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}

		if doesTitleWithTheSameNameExist {
			c.AbortWithStatusJSON(409, gin.H{"error": "тайтл с таким названием уже существует"})
			return
		}

		editedTitle.Name = sql.NullString{String: name, Valid: true}
	}

	err = tx.Raw(
		`INSERT INTO titles_on_moderation (created_at, name, description, author_id, creator_id, existing_id)
		VALUES (NOW(), ?, ?, ?, ?, ?)
		ON CONFLICT (existing_id) DO UPDATE
		SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			author_id = EXCLUDED.author_id,
			creator_id = EXCLUDED.creator_id,
			updated_at = NOW()
		RETURNING id`,
		editedTitle.Name, editedTitle.Description, editedTitle.AuthorID, editedTitle.CreatorID, editedTitle.ExistingID,
	).Scan(&editedTitle.ID).Error

	if err != nil {
		if dbErrors.IsForeignKeyViolation(err, constraints.FkTitlesOnModerationAuthor) {
			c.AbortWithStatusJSON(404, gin.H{"error": "автор не найден"})
			return
		}

		if dbErrors.IsUniqueViolation(err, constraints.UniTitlesOnModerationName) {
			c.AbortWithStatusJSON(409, gin.H{"error": "тайтл с таким названием уже ожидает модерации"})
			return
		}

		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(genresIDs) != 0 {
		err = tx.Exec(
			`INSERT INTO title_on_moderation_genres (title_on_moderation_id, genre_id)
			SELECT ?, UNNEST(?::BIGINT[])`,
			editedTitle.ID, pq.Array(form.Value["genresIds"]),
		).Error

		if err != nil {
			if dbErrors.IsForeignKeyViolation(err, constraints.FkTitleOnModerationGenresGenre) {
				c.AbortWithStatusJSON(409, gin.H{"error": "жанры не найдены"})
			} else {
				log.Println(err)
				c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			}
			return
		}
	}

	if coverFileHeader != nil {
		cover, err := utils.ReadMultipartFile(coverFileHeader, 2<<20)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}

		filter := bson.M{"title_on_moderation_id": editedTitle.ID}
		update := bson.M{"$set": bson.M{"cover": cover}}
		opts := options.Update().SetUpsert(true)

		if _, err := h.TitlesOnModerationCovers.UpdateOne(c.Request.Context(), filter, update, opts); err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "изменения тайтла успешно отправлены на модерацию"})

	if _, err := h.NotificationsClient.NotifyAboutTitleOnModeration(c.Request.Context(), &pb.TitleOnModeration{ID: uint64(*editedTitle.ExistingID), New: false}); err != nil {
		log.Println(err)
	}
}

func parseEditTitleParams(form *multipart.Form, paramFn func(string) string) (name, description string, titleID, authorID uint, genresIDs []uint, coverFileHeader *multipart.FileHeader, err error) {
	titleIDuint64, err := strconv.ParseUint(paramFn("id"), 10, 64)
	if err != nil {
		return "", "", 0, 0, nil, nil, errors.New("указан невалидный id тайтла")
	}

	if len(form.Value["name"]) == 0 && len(form.Value["description"]) == 0 && len(form.Value["genresIds"]) == 0 && len(form.File["cover"]) == 0 && len(form.Value["authorId"]) == 0 {
		return "", "", 0, 0, nil, nil, errors.New("запрос должен содержать как минимум 1 изменяемый параметр")
	}

	if len(form.Value["name"]) != 0 {
		name = form.Value["name"][0]
	}

	if len(form.Value["description"]) != 0 {
		description = form.Value["description"][0]
	}

	if len(form.Value["authorId"]) != 0 {
		authorIDuint64, err := strconv.ParseUint(form.Value["authorId"][0], 10, 64)
		if err != nil {
			return "", "", 0, 0, nil, nil, errors.New("указан невалидный id автора")
		}
		authorID = uint(authorIDuint64)
	}

	if len(form.Value["genresIds"]) != 0 {
		genresIDs, err = utils.ParseUintSlice(form.Value["genresIds"])
		if err != nil {
			return "", "", 0, 0, nil, nil, errors.New("указаны невалидные id жанров")
		}
	}

	if len(form.File["cover"]) != 0 {
		coverFileHeader = form.File["cover"][0]
		if coverFileHeader.Size > 2<<20 {
			return "", "", 0, 0, nil, nil, errors.New("превышен максимальный размер обложки (2мб)")
		}
	}

	return name, description, uint(titleIDuint64), authorID, genresIDs, coverFileHeader, nil
}
