package moderation

import (
	"fmt"
	"log"
	"strings"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
)

type getMyAuthorsOnModerationParams struct {
	dto.CommonParams
}

func (h handler) GetMyAuthorsOnModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var params getMyAuthorsOnModerationParams

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

	selects.WriteString("id, name, original_name, english_name")

	if params.Query != nil {
		selects.WriteString(",name <-> ? AS name_distance, english_name <-> ? AS english_name_distance, original_name <-> ? AS original_name_distance")
		args = append(args, *params.Query, *params.Query, *params.Query)
	}

	query := h.DB.Table("authors_on_moderation").
		Select(selects.String(), args...).
		Where("creator_id = ?", claims.ID).
		Limit(int(params.Limit)).
		Offset(offset)

	if params.Query != nil {
		query = query.
			Where("name % ? OR english_name % ? OR original_name % ?", *params.Query, *params.Query, *params.Query).
			Order("name_distance, english_name_distance, original_name_distance")
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
		c.AbortWithStatusJSON(404, gin.H{"error": "не найдено ваших авторов на модерации"})
		return
	}

	c.JSON(200, &result)
}
