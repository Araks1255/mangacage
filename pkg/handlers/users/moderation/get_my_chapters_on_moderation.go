package moderation

import (
	"fmt"
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
)

type getMyChaptersOnModerationParams struct {
	Sort  string  `form:"sort"`
	Query *string `form:"query"`
	Order string  `form:"order"`
	Page  int     `form:"page,default=1"`
	Limit uint    `form:"limit,default=20"`

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

	query := h.DB.Table("chapters_on_moderation AS com").
		Select(
			`com.*, c.name AS existing,
			t.name AS title, tom.name AS title_on_moderation, teams.name AS team`,
		).
		Joins("LEFT JOIN chapters AS c ON com.existing_id = c.id").
		Joins("LEFT JOIN titles AS t ON com.title_id = t.id").
		Joins("INNER JOIN teams ON teams.id = com.team_id").
		Joins("LEFT JOIN titles_on_moderation AS tom ON com.title_on_moderation_id = tom.id").
		Where("com.creator_id = ?", claims.ID).
		Offset(offset).
		Limit(int(params.Limit))

	if params.Query != nil {
		query = query.Where("lower(com.name) ILIKE lower(?)", fmt.Sprintf("%%%s%%", *params.Query))
	}

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
		query = query.Where("t.id = ?", params.TitleID)
	}
	if params.TitleOnModerationID != nil {
		query = query.Where("tom.id = ?", params.TitleOnModerationID)
	}

	if params.Order != "desc" && params.Order != "asc" {
		params.Order = "desc"
	}

	switch params.Sort {
	case "createdAt":
		query = query.Order(fmt.Sprintf("com.id %s", params.Order))
	case "numberOfPages":
		query = query.Order(fmt.Sprintf("com.number_of_pages %s", params.Order))
	default:
		query = query.Order(fmt.Sprintf("com.name %s", params.Order))
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
