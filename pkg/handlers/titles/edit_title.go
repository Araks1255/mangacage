package titles

import (
	"context"
	"errors"
	"log"
	"mime/multipart"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/common/utils"
	"github.com/Araks1255/mangacage/pkg/handlers/helpers"
	"github.com/Araks1255/mangacage/pkg/handlers/helpers/titles"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm/clause"
)

func (h handler) EditTitle(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var requestBody models.TitleOnModerationDTO
	c.ShouldBindWith(&requestBody, binding.FormMultipart)

	if requestBody.Cover != nil && requestBody.Cover.Size > 2<<20 {
		c.AbortWithStatusJSON(400, gin.H{"error": "превышен максимальный размер обложки (2мб)"})
		return
	}

	editedTitle, code, err := mapEditTitleParamsToTitleOnModeration(claims.ID, &requestBody, c.Param)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	{
		ok, err := titles.IsUserTeamTranslatingTitle(tx, claims.ID, *editedTitle.ExistingID)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}

		if !ok {
			c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден среди тайтлов, переводимых вашей командой"})
			return
		}
	}

	{
		if editedTitle.Name != nil || editedTitle.EnglishName != nil || editedTitle.OriginalName != nil {
			exists, err := helpers.CheckEntityWithTheSameNameExistence(tx, "titles", *editedTitle.Name, editedTitle.EnglishName, editedTitle.OriginalName)
			if err != nil {
				log.Println(err)
				c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
				return
			}

			if exists {
				c.AbortWithStatusJSON(409, gin.H{"error": "тайтл с таким названием уже существует"})
				return
			}
		}
	}

	{
		onConflictClause := clause.OnConflict{ // Тут то же самое, что было в сыром SQL.
			Columns:   []clause.Column{{Name: "existing_id"}}, // Конфликт по existing_id
			UpdateAll: true,                                   // Столбцы обновляются по структуре
		}

		err = tx.Clauses(onConflictClause).Create(&editedTitle).Error

		if err != nil {
			code, err := titles.ParseTitleOnModerationInsertError(err)
			if code == 500 {
				log.Println(err)
			}
			c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
	}

	{
		if len(requestBody.GenresIDs) != 0 {
			code, err := titles.InsertTitleOnModerationGenres(tx, editedTitle.ID, requestBody.GenresIDs)
			if err != nil {
				if code == 500 {
					log.Println(err)
				}
				c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
				return
			}
		}
	}

	{
		if len(requestBody.TagsIDs) != 0 {
			code, err := titles.InsertTitleOnModerationTags(tx, editedTitle.ID, requestBody.TagsIDs)
			if err != nil {
				if code == 500 {
					log.Println(err)
				}
				c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
				return
			}
		}
	}

	{
		if requestBody.Cover != nil {
			if err := upsertTitleOnModerationCover(c.Request.Context(), h.TitlesOnModerationCovers, requestBody.Cover, editedTitle.ID); err != nil {
				log.Println(err)
				c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
				return
			}
		}
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "изменения тайтла успешно отправлены на модерацию"})

	if _, err := h.NotificationsClient.NotifyAboutTitleOnModeration(c.Request.Context(), &pb.TitleOnModeration{ID: uint64(*editedTitle.ExistingID), New: false}); err != nil {
		log.Println(err)
	}
}

func mapEditTitleParamsToTitleOnModeration(userID uint, body *models.TitleOnModerationDTO, paramFn func(string) string) (res *models.TitleOnModeration, code int, err error) {
	titleID, err := strconv.ParseUint(paramFn("id"), 10, 64)
	if err != nil {
		return nil, 400, errors.New("указан невалидный id тайтла")
	}

	ok, err := utils.HasAnyNonEmptyFields(body)
	if err != nil {
		return nil, 500, err
	}

	if !ok {
		return nil, 400, errors.New("запрос должен содержать как минимум 1 изменяемый параметр")
	}

	titleIDuint := uint(titleID)

	titleOnModeration := body.ToTitleOnModeration(userID, &titleIDuint)

	return &titleOnModeration, 0, nil
}

func upsertTitleOnModerationCover(ctx context.Context, collection *mongo.Collection, coverFileHeader *multipart.FileHeader, titleOnModerationID uint) error {
	cover, err := utils.ReadMultipartFile(coverFileHeader, 2<<20)
	if err != nil {
		return err
	}

	filter := bson.M{"title_on_moderation_id": titleOnModerationID}
	update := bson.M{"$set": bson.M{"cover": cover}}
	opts := options.Update().SetUpsert(true)

	if _, err := collection.UpdateOne(ctx, filter, update, opts); err != nil {
		return err
	}

	return nil
}
