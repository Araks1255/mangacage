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
	"gorm.io/gorm"
)

func (h handler) EditTitle(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var requestBody dto.EditTitleDTO

	if err := c.ShouldBindWith(&requestBody, binding.FormMultipart); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
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

	code, err = checkEditTitleConflicts(h.DB, *editedTitle, claims.ID)
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
		code, err := titles.UpsertTitleOnModerationGenres(tx, editedTitle.ID, requestBody.GenresIDs)
		if err != nil {
			if code == 500 {
				log.Println(err)
			}
			c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
	}

	if len(requestBody.TagsIDs) != 0 {
		code, err := titles.UpsertTitleOnModerationTags(tx, editedTitle.ID, requestBody.TagsIDs)
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
	if body.Cover != nil && body.Cover.Size > 2<<20 {
		return nil, 400, errors.New("превышен максимальный размер обложки (2мб)")
	}

	titleID, err := strconv.ParseUint(paramFn("id"), 10, 64)
	if err != nil {
		return nil, 400, errors.New("указан невалидный id тайтла")
	}

	ok, err := utils.HasAnyNonEmptyFields(body, "AuthorID")
	if err != nil {
		return nil, 500, err
	}

	if !ok {
		return nil, 400, errors.New("запрос должен содержать как минимум 1 изменяемый параметр")
	}

	titleOnModeration := body.ToTitleOnModeration(userID, uint(titleID))

	return &titleOnModeration, 0, nil
}

func checkEditTitleConflicts(db *gorm.DB, title models.TitleOnModeration, userID uint) (code int, err error) {
	ok, err := titles.IsUserTeamTranslatingTitle(db, userID, *title.ExistingID)
	if err != nil {
		return 500, err
	}

	if !ok {
		return 404, errors.New("тайтл не найден среди тайтлов, переводимых вашей командой")
	}

	if title.Name != nil || title.EnglishName != nil || title.OriginalName != nil {
		exists, err := helpers.CheckEntityWithTheSameNameExistence(db, "titles", title.Name, title.EnglishName, title.OriginalName)
		if err != nil {
			return 500, err
		}

		if exists {
			return 409, errors.New("тайтл с таким названием уже существует")
		}
	}

	return 0, nil
}
