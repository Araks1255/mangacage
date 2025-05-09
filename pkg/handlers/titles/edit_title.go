package titles

import (
	"database/sql"
	"log"
	"slices"
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

	titleID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id тайтла должен быть числом"})
		return
	}

	var userRoles []string
	h.DB.Raw(
		`SELECT r.name FROM roles AS r
		INNER JOIN user_roles AS ur ON r.id = ur.role_id
		WHERE ur.user_id = ?`, claims.ID,
	).Scan(&userRoles)

	if !slices.Contains(userRoles, "team_leader") && !slices.Contains(userRoles, "ex_team_leader") && !slices.Contains(userRoles, "moder") && !slices.Contains(userRoles, "admin") {
		c.AbortWithStatusJSON(403, gin.H{"error": "вы не являетесь лидером команды перевода"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	if len(form.Value["name"]) == 0 && len(form.Value["description"]) == 0 && len(form.Value["genresIds"]) == 0 && len(form.File["cover"]) == 0 && len(form.Value["authorId"]) == 0 {
		c.AbortWithStatusJSON(400, gin.H{"error": "запрос должен содержать хотя-бы один изменяемый параметр"})
		return
	}

	if len(form.File["cover"]) != 0 && form.File["cover"][0].Size > 10<<20 {
		c.AbortWithStatusJSON(400, gin.H{"error": "превышен лимит размера обложки (10мб)"})
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
		ExistingID: sql.NullInt64{Int64: int64(titleID), Valid: true},
		CreatorID:  claims.ID,
	}

	if len(form.Value["authorId"]) != 0 {
		desiredAuthorID, err := strconv.ParseUint(form.Value["authorId"][0], 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id автора"})
			return
		}

		editedTitle.AuthorID = sql.NullInt64{Int64: int64(desiredAuthorID), Valid: true}
	}

	if len(form.Value["name"]) != 0 {
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

		editedTitle.Name = sql.NullString{String: form.Value["name"][0], Valid: true}
	}

	if len(form.Value["description"]) != 0 {
		editedTitle.Description = form.Value["description"][0]
	}

	err = tx.Raw(
		`INSERT INTO titles_on_moderation (created_at, name, description, author_id, creator_id, existing_id)
		VALUES (NOW(), ?, ?, ?, ?, ?, ?)
		ON CONFLICT (existing_id) DO UPDATE
		SET
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			author_id = EXCLUDED.author_id,
			genres = EXCLUDED.genres,
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
	}

	if len(form.Value["genresIds"]) != 0 {
		err = tx.Exec(
			`INSERT INTO titles_on_moderation_genres (title_id, genre_id)
			SELECT ?, UNNEST(?::INTEGER[])`,
			pq.Array(form.Value["genresIds"]),
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
	}

	if len(form.File["cover"]) != 0 {
		cover, err := utils.ReadMultipartFile(form.File["cover"][0], 10<<20)
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

	if _, err := h.NotificationsClient.NotifyAboutTitleOnModeration(c.Request.Context(), &pb.TitleOnModeration{ID: uint64(editedTitle.ExistingID.Int64), New: false}); err != nil {
		log.Println(err)
	}
}
