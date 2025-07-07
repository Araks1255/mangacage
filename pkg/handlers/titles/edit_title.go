package titles

import (
	"errors"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/Araks1255/mangacage/pkg/common/utils"
	"github.com/Araks1255/mangacage/pkg/handlers/helpers"
	"github.com/Araks1255/mangacage/pkg/handlers/helpers/titles"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func (h handler) EditTitle(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var requestBody dto.EditTitleDTO

	if err := c.ShouldBindWith(&requestBody, binding.FormMultipart); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

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

	if editedTitle.Name != nil || editedTitle.EnglishName != nil || editedTitle.OriginalName != nil {
		exists, err := helpers.CheckEntityWithTheSameNameExistence(tx, "titles", editedTitle.Name, editedTitle.EnglishName, editedTitle.OriginalName)
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

	err = tx.Clauses(helpers.OnExistingIDConflictClause).Create(&editedTitle).Error

	if err != nil {
		code, err := titles.ParseTitleOnModerationInsertError(err)
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

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

	if requestBody.Cover != nil {
		if err := titles.UpsertTitleOnModerationCover(c.Request.Context(), h.TitlesCovers, requestBody.Cover, editedTitle.ID, claims.ID); err != nil {
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

func mapEditTitleParamsToTitleOnModeration(userID uint, body *dto.EditTitleDTO, paramFn func(string) string) (res *models.TitleOnModeration, code int, err error) {
	titleID, err := strconv.ParseUint(paramFn("id"), 10, 64)
	if err != nil {
		return nil, 400, errors.New("указан невалидный id тайтла")
	}

	ok, err := utils.HasAnyNonEmptyFields(body, "AuthorID", "AuthorOnModerationID")
	if err != nil {
		return nil, 500, err
	}

	if !ok {
		return nil, 400, errors.New("запрос должен содержать как минимум 1 изменяемый параметр")
	}

	titleOnModeration := body.ToTitleOnModeration(userID, uint(titleID))

	return &titleOnModeration, 0, nil
}
