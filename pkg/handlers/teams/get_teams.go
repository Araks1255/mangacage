package teams

import (
	"fmt"
	"log"
	"strings"

	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
)

type getTeamParams struct {
	dto.CommonParams
	TranslatingTitleID       *uint `form:"translatingTitleId"`
	NumberOfParticipantsFrom *uint `form:"numberOfParticipantsFrom"`
	NumberOfParticipantsTo   *uint `form:"numberOfParticipantsTo"`
}

func (h handler) GetTeams(c *gin.Context) {
	var params getTeamParams

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

	selects.WriteString("t.id, t.name")

	if params.Query != nil {
		selects.WriteString(",t.name <-> ? AS distance")
		args = append(args, *params.Query)
	}

	query := h.DB.Table("teams AS t").Select(selects.String(), args...).Limit(int(params.Limit)).Offset(offset)

	if params.NumberOfParticipantsFrom != nil {
		query = query.Where("t.number_of_participants >= ?", params.NumberOfParticipantsFrom)
	}
	if params.NumberOfParticipantsTo != nil {
		query = query.Where("t.number_of_participants <= ?", params.NumberOfParticipantsTo)
	}

	if params.TranslatingTitleID != nil {
		query = query.Joins("INNER JOIN title_teams AS tt ON tt.team_id = t.id").
			Where("tt.title_id = ?", *params.TranslatingTitleID)
	}

	if params.Query != nil {
		query = query.Where("t.name % ?", *params.Query).Order("distance ASC")
	} else {
		if params.Order != "desc" && params.Order != "asc" {
			params.Order = "asc"
		}

		switch params.Sort {
		case "createdAt":
			query = query.Order(fmt.Sprintf("t.id %s", params.Order))
		default:
			query = query.Order(fmt.Sprintf("t.name %s", params.Order))
		}
	}

	var result []dto.ResponseTeamDTO

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
