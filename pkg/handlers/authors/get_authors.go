package authors

import (
	"fmt"
	"log"
	"strings"

	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
)

type GetAuthorsParams struct {
	dto.CommonParams
}

func (h handler) GetAuthors(c *gin.Context) {
	var params GetAuthorsParams

	if err := c.ShouldBindQuery(&params); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	offset := (params.Page - 1) * int(params.Limit)
	if offset < 0 {
		offset = 0
	}

	var selects strings.Builder
	args := make([]any, 0, 3)

	selects.WriteString("id, name, english_name, original_name")

	if params.Query != nil {
		selects.WriteString(",name <-> ? AS name_distance, english_name <-> ? AS english_name_distance, original_name <-> ? AS original_name_distance")
		args = append(args, *params.Query, *params.Query, *params.Query)
	}

	query := h.DB.Table("authors").Select(selects.String(), args...).Limit(int(params.Limit)).Offset(offset)

	if params.Query != nil {
		query = query.
			Where("name % ? OR english_name % ? OR original_name % ?", *params.Query, *params.Query, *params.Query).
			Order("name_distance, english_name_distance, original_name_distance ASC")
	} else {
		if params.Order != "desc" && params.Order != "asc" {
			params.Order = "asc"
		}

		switch params.Sort {
		case "createdAt":
			query = query.Order(fmt.Sprintf("id %s", params.Order))
		default:
			query = query.Order(fmt.Sprintf("name %s", params.Order))
		}
	}

	var result []dto.ResponseAuthorDTO

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
