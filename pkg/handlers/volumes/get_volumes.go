package volumes

import (
	"fmt"
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
)

type GetVolumesParams struct {
	Sort  string  `form:"sort"`
	Query *string `form:"query"`
	Order string  `form:"order"`
	Page  int     `form:"page,default=1"`
	Limit uint    `form:"limit,default=20"`

	TitleID *uint `form:"titleId"`
	TeamID  *uint `form:"teamId"`
}

func (h handler) GetVolumes(c *gin.Context) {
	var params GetVolumesParams

	if err := c.ShouldBindQuery(&params); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	query := h.DB.Table("volumes AS v").
		Select("v.*, t.name AS title, teams.name AS team").
		Joins("INNER JOIN titles AS t ON t.id = v.title_id").
		Joins("INNER JOIN teams ON teams.id = v.team_id")

	if params.Query != nil {
		query = query.Where("lower(v.name) ILIKE lower(?)", fmt.Sprintf("%%%s%%", *params.Query))
	}

	if params.TeamID != nil {
		query = query.Where("v.team_id = ?", params.TeamID)
	}

	if params.TitleID != nil {
		query = query.Where("t.id = ?", params.TitleID)
	}

	if params.Order != "desc" && params.Order != "asc" {
		params.Order = "desc"
	}

	switch params.Sort {
	case "createdAt":
		query = query.Order(fmt.Sprintf("v.id %s", params.Order))

	default:
		query = query.Order(fmt.Sprintf("v.name %s", params.Order))
	}

	offset := (params.Page - 1) * int(params.Limit)
	if offset < 0 {
		offset = 0
	}

	query = query.Limit(int(params.Limit)).Offset(offset)

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
