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
	"github.com/Araks1255/mangacage_protos/gen/enums"
	pb "github.com/Araks1255/mangacage_protos/gen/site_notifications"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

func (h handler) EditTitle(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var requestBody dto.EditTitleDTO

	if err := c.ShouldBindWith(&requestBody, binding.FormMultipart); err != nil {
		log.Println(requestBody)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	editedTitle, code, err := mapEditTitleRequestBodyToTitleOnModeration(claims.ID, &requestBody, c.Param)
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

	err = helpers.UpsertEntityChanges(tx, editedTitle, *editedTitle.ExistingID)
	if err != nil {
		code, err := titles.ParseTitleOnModerationInsertError(err)
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	if requestBody.Cover != nil {
		if code, err := titles.CreateTitleOnModerationCover(tx, h.PathToMediaDir, editedTitle.ID, requestBody.Cover); err != nil {
			if code == 500 {
				log.Println(err)
			}
			c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
	}

	if code, err = titles.UpsertTitleOnModerationGenres(tx, editedTitle.ID, requestBody.GenresIDs); err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	if code, err = titles.UpsertTitleOnModerationTags(tx, editedTitle.ID, requestBody.TagsIDs); err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "изменения тайтла успешно отправлены на модерацию"})

	if _, err := h.NotificationsClient.NotifyAboutNewModerationRequest(
		c.Request.Context(),
		&pb.ModerationRequest{
			EntityOnModeration: enums.EntityOnModeration_ENTITY_ON_MODERATION_TITLE,
			ID:                 uint64(editedTitle.ID),
		},
	); err != nil {
		log.Println(err)
	}
}

func mapEditTitleRequestBodyToTitleOnModeration(
	userID uint,
	body *dto.EditTitleDTO,
	paramFn func(string) string,
) (
	res *models.TitleOnModeration,
	code int,
	err error,
) {
	if body.Cover != nil && body.Cover.Size > 2<<20 {
		return nil, 400, errors.New("превышен максимальный размер обложки (2мб)")
	}

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

	return body.ToTitleOnModeration(userID, uint(titleID)), 0, nil
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
