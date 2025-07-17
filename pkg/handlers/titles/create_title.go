package titles

import (
	"errors"
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"gorm.io/gorm"

	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/Araks1255/mangacage/pkg/handlers/helpers"
	"github.com/Araks1255/mangacage/pkg/handlers/helpers/titles"
	pb "github.com/Araks1255/mangacage_protos"
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

	code, err := checkCreateTitleConflicts(h.DB, requestBody, claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	newTitle := requestBody.ToTitleOnModeration(claims.ID)

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	err = tx.Clauses(helpers.OnIDConflictClause).Create(&newTitle).Error

	if err != nil {
		code, err := titles.ParseTitleOnModerationInsertError(err)
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	code, err = titles.UpsertTitleOnModerationGenres(tx, newTitle.ID, requestBody.GenresIDs)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	code, err = titles.UpsertTitleOnModerationTags(tx, newTitle.ID, requestBody.TagsIDs)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	if err := titles.UpsertTitleOnModerationCover(c.Request.Context(), h.TitlesCovers, requestBody.Cover, newTitle.ID, claims.ID); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "тайтл успешно отправлен на модерацию"})

	if _, err := h.NotificationsClient.NotifyAboutTitleOnModeration(c.Request.Context(), &pb.TitleOnModeration{ID: uint64(newTitle.ID), New: true}); err != nil {
		log.Println(err)
	}
}

func checkCreateTitleConflicts(db *gorm.DB, requestBody dto.CreateTitleDTO, userID uint) (code int, err error) {
	if requestBody.AuthorID != nil && requestBody.AuthorOnModerationID != nil {
		return 400, errors.New("должен быть заполнен только один id автора")
	}
	if requestBody.Cover.Size > 2<<20 {
		return 400, errors.New("превышен максимальный размер обложки (2мб)")
	}

	exists, err := helpers.CheckEntityWithTheSameNameExistence(db, "titles", &requestBody.Name, &requestBody.EnglishName, &requestBody.OriginalName)
	if err != nil {
		return 500, err
	}

	if exists {
		return 409, errors.New("тайтл с таким названием уже существует")
	}

	if requestBody.ID != nil {
		isOwner, err := helpers.CheckEntityOnModerationOwnership(db, "titles", *requestBody.ID, userID)
		if err != nil {
			return 500, err
		}

		if !isOwner {
			return 403, errors.New("изменять заявку на модерацию может только её создатель")
		}
	}

	return 0, nil
}
