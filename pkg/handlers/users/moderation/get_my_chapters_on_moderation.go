package moderation

import (
	"fmt"
	"log"
	"strings"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
)

type getMyChaptersOnModerationParams struct {
	dto.CommonParams

	NumberOfPagesFrom *int `form:"numberOfPagesFrom"`
	NumberOfPagesTo   *int `form:"numberOfPagesTo"`

	Volume *uint `form:"volume"`

	TitleID             *uint `form:"titleId"`
	TitleOnModerationID *uint `form:"titleOnModerationId"`

	ModerationType string `form:"type"`
}

func (h handler) GetMyChaptersOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var params getMyChaptersOnModerationParams

	if err := c.ShouldBindQuery(&params); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	offset := (params.Page - 1) * int(params.Limit)
	if offset < 0 {
		offset = 0
	}

	var selects strings.Builder
	args := make([]any, 0, 1)

	selects.WriteString("com.id, com.name, com.title_id, com.title_on_moderation_id, com.team_id, teams.name AS team")

	if params.Query != nil {
		selects.WriteString(",com.name <-> ? AS distance")
		args = append(args, *params.Query)
	}

	query := h.DB.Table("chapters_on_moderation AS com").
		Select(selects.String(), args...).
		Where("com.creator_id = ?", claims.ID).
		Joins("LEFT JOIN teams ON com.team_id = teams.id").
		Offset(offset).
		Limit(int(params.Limit))

	if params.ModerationType == "new" {
		query = query.Where("com.existing_id IS NULL")
	}
	if params.ModerationType == "edited" {
		query = query.Where("com.existing_id IS NOT NULL")
	}

	if params.NumberOfPagesFrom != nil {
		query = query.Where("com.number_of_pages >= ?", params.NumberOfPagesFrom)
	}
	if params.NumberOfPagesTo != nil {
		query = query.Where("com.number_of_pages <= ?", params.NumberOfPagesTo)
	}

	if params.Volume != nil {
		query = query.Where("com.volume = ?", params.Volume)
	}

	if params.TitleID != nil {
		query = query.Where("com.title_id = ?", params.TitleID)
	}
	if params.TitleOnModerationID != nil {
		query = query.Where("com.title_on_moderation_id = ?", params.TitleOnModerationID)
	}

	if params.Query != nil {
		query = query.Where("com.name % ?", *params.Query).Order("distance ASC")
	} else {
		if params.Order != "desc" && params.Order != "asc" {
			params.Order = "asc"
		}

		switch params.Sort {
		case "createdAt":
			query = query.Order(fmt.Sprintf("com.id %s", params.Order))
		case "numberOfPages":
			query = query.Order(fmt.Sprintf("com.number_of_pages %s", params.Order))
		case "number":
			query = query.Order(fmt.Sprintf("CAST(substring(com.name from '[0-9.]+') AS DECIMAL) %s", params.Order))
		default:
			query = query.Order(fmt.Sprintf("com.name %s", params.Order))
		}
	}

	var result []dto.ResponseChapterDTO

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
