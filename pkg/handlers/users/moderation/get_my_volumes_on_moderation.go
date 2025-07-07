package moderation

import (
	"fmt"
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
)

type getMyVolumesOnModerationParams struct {
	Sort  string  `form:"sort"`
	Query *string `form:"query"`
	Order string  `form:"order"`
	Page  int     `form:"page,default=1"`
	Limit uint    `form:"limit,default=20"`

	ModerationType      string `form:"type"`
	TitleID             *uint  `form:"titleId"`
	TitleOnModerationID *uint  `form:"titleOnModerationId"`
}

func (h handler) GetMyVolumesOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var params getMyVolumesOnModerationParams

	if err := c.ShouldBindQuery(&params); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	offset := (params.Page - 1) * int(params.Limit)
	if offset < 0 {
		offset = 0
	}

	query := h.DB.Table("volumes_on_moderation AS vom").
		Select("vom.*, t.name AS title, tom.name AS title_on_moderation, v.name AS existing").
		Joins("LEFT JOIN volumes AS v ON vom.existing_id = v.id").
		Joins("LEFT JOIN titles AS t ON vom.title_id = t.id").
		Joins("LEFT JOIN titles_on_moderation AS tom ON vom.title_on_moderation_id = tom.id").
		Where("vom.creator_id = ?", claims.ID).
		Offset(offset).Limit(int(params.Limit))

	if params.ModerationType == "new" {
		query = query.Where("vom.existing_id IS NULL")
	}
	if params.ModerationType == "edited" {
		query = query.Where("vom.existing_id IS NOT NULL")
	}

	if params.TitleID != nil {
		query = query.Where("vom.title_id = ?", *params.TitleID)
	}
	if params.TitleOnModerationID != nil {
		query = query.Where("vom.title_on_moderation_id = ?", *params.TitleOnModerationID)
	}

	if params.Order != "desc" && params.Order != "asc" {
		params.Order = "desc"
	}

	switch params.Sort {
	case "createdAt":
		query = query.Order(fmt.Sprintf("vom.id %s", params.Order))
	default:
		query = query.Order(fmt.Sprintf("vom.name %s", params.Order))
	}

	var result []dto.ResponseVolumeDTO

	if err := query.Scan(&result).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if len(result) == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "по вашему запросу ничего не найдено"})
		return
	}

	c.JSON(200, &result)
}
