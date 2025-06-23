package authors

import (
	"fmt"
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
)

type GetAuthorsParams struct {
	Query *string `form:"query"`
	Order string  `form:"order"`
	Sort  string  `form:"sort"`
	Page  int     `form:"page,default=1"`
	Limit uint    `form:"limit,default=20"`
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

	query := h.DB.Table("authors").Select("*").Limit(int(params.Limit)).Offset(offset)

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

	var result []models.AuthorDTO

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
