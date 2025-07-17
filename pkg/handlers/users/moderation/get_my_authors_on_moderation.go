package moderation

import (
	"fmt"
	"log"

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

	query := h.DB.Table("authors_on_moderation").
		Select("*").
		Where("creator_id = ?", claims.ID).
		Limit(int(params.Limit)).
		Offset(offset)

	if params.Query != nil {
		query = query.Where("lower(name) ILIKE lower(?)", fmt.Sprintf("%%%s%%", *params.Query))
	}

	if params.Order != "desc" && params.Order != "asc" {
		params.Order = "desc"
	}

	switch params.Sort {
	case "createdAt":
		query = query.Order(fmt.Sprintf("id %s", params.Order))
	default:
		query = query.Order(fmt.Sprintf("name %s", params.Order))
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
