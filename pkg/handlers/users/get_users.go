package users

import (
	"fmt"
	"log"
	"strings"

	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
)

type getUsersParams struct {
	dto.CommonParams

	TeamID *uint `form:"teamId"`
}

func (h handler) GetUsers(c *gin.Context) {
	var params getUsersParams

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

	selects.WriteString(
		`u.*, ARRAY(
			SELECT r.name FROM roles AS r
			INNER JOIN user_roles AS ur ON ur.role_id = r.id
			WHERE ur.user_id = u.id AND r.type = 'team'
		) AS roles`,
	)

	if params.Query != nil {
		selects.WriteString(",u.user_name <-> ? AS distance")
		args = append(args, *params.Query)
	}

	query := h.DB.Table("users AS u").
		Select(selects.String(), args...).
		Limit(int(params.Limit)).Offset(offset)

	if params.TeamID != nil {
		query = query.Where("u.team_id = ?", params.TeamID)
	} else {
		query = query.Where("visible")
	}

	if params.Order != "desc" && params.Order != "asc" {
		params.Order = "desc"
	}

	if params.Query != nil {
		query = query.Where("user_name % ?", *params.Query).Order("distance ASC")
	} else {
		switch params.Sort {
		case "createdAt":
			query = query.Order(fmt.Sprintf("u.id %s", params.Order))
		default:
			query = query.Order(fmt.Sprintf("u.user_name %s", params.Order))
		}
	}

	var result []dto.ResponseUserDTO

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
