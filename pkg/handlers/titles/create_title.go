package titles

import (
	"errors"
	_ "image/jpeg"
	_ "image/png"
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"gorm.io/gorm"

	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/Araks1255/mangacage/pkg/handlers/helpers"
	"github.com/Araks1255/mangacage/pkg/handlers/helpers/titles"
	"github.com/Araks1255/mangacage_protos/gen/enums"
	pb "github.com/Araks1255/mangacage_protos/gen/site_notifications"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func (h handler) CreateTitle(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var requestBody dto.CreateTitleDTO

	if err := c.ShouldBindWith(&requestBody, binding.FormMultipart); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	if code, err := checkCreateTitleConflicts(h.DB, &requestBody, claims.ID); err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	newTitle := requestBody.ToTitleOnModeration(claims.ID)

	err := helpers.UpsertEntityOnModeration(tx, newTitle, newTitle.ID)

	if err != nil {
		code, err := titles.ParseTitleOnModerationInsertError(err)
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	if code, err := titles.CreateTitleOnModerationCover(tx, h.PathToMediaDir, newTitle.ID, requestBody.Cover); err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	if code, err := titles.UpsertTitleOnModerationGenres(tx, newTitle.ID, requestBody.GenresIDs); err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	if code, err := titles.UpsertTitleOnModerationTags(tx, newTitle.ID, requestBody.TagsIDs); err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "тайтл успешно отправлен на модерацию", "id": newTitle.ID})

	if _, err := h.NotificationsClient.NotifyAboutNewModerationRequest(
		c.Request.Context(),
		&pb.ModerationRequest{
			EntityOnModeration: enums.EntityOnModeration_ENTITY_ON_MODERATION_TITLE,
			ID:                 uint64(newTitle.ID),
		},
	); err != nil {
		log.Println(err)
	}
}

func checkCreateTitleConflicts(db *gorm.DB, requestBody *dto.CreateTitleDTO, userID uint) (code int, err error) {
	if requestBody.AuthorID != nil && requestBody.AuthorOnModerationID != nil {
		return 400, errors.New("должен быть заполнен только один id автора")
	}
	if requestBody.Cover.Size > 2<<20 {
		return 400, errors.New("превышен максимальный размер обложки (2мб)")
	}

	if requestBody.ID != nil {
		isOwner, err := helpers.CheckEntityOnModerationOwnership(db, "titles", *requestBody.ID, userID)
		if err != nil {
			return 500, err
		}

		if !isOwner {
			return 403, errors.New("изменять заявку на модерацию может только её создатель")
		}

		return 0, nil
	}

	exists, err := helpers.CheckEntityWithTheSameNameExistence(db, "titles", &requestBody.Name, &requestBody.EnglishName, &requestBody.OriginalName)

	if err != nil {
		return 500, err
	}

	if exists {
		return 409, errors.New("тайтл с таким названием уже существует")
	}

	return 0, nil
}
