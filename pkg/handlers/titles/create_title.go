package titles

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"

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

	if requestBody.AuthorID != nil && requestBody.AuthorOnModerationID != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "должен быть заполнен только один id автора"})
		return
	}

	if requestBody.Cover.Size > 2<<20 {
		c.AbortWithStatusJSON(400, gin.H{"error": "превышен максимальный размер обложки (2мб)"})
		return
	}

	newTitle := requestBody.ToTitleOnModeration(claims.ID)

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	exists, err := helpers.CheckEntityWithTheSameNameExistence(tx, "titles", &requestBody.Name, &requestBody.EnglishName, &requestBody.OriginalName)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if exists {
		c.AbortWithStatusJSON(409, gin.H{"error": "тайтл с таким названием уже существует"})
		return
	}

	err = tx.Clauses(helpers.OnIDConflictClause).Create(&newTitle).Error

	if err != nil {
		code, err := titles.ParseTitleOnModerationInsertError(err)
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	code, err := titles.InsertTitleOnModerationGenres(tx, newTitle.ID, requestBody.GenresIDs)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	code, err = titles.InsertTitleOnModerationTags(tx, newTitle.ID, requestBody.TagsIDs)
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
