package genres

import (
	"fmt"
	"log"
	"strings"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
)

type getGenresParams struct {
	dto.CommonParams

	MyFavorites *bool `form:"myFavorites" binding:"excluded_with=FavoritedBy"`
}

func (h handler) GetGenres(c *gin.Context) {
	var params getGenresParams

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

	selects.WriteString("g.id, g.name")

	if params.Query != nil {
		selects.WriteString(",g.name <-> ? AS distance")
		args = append(args, *params.Query)
	}

	query := h.DB.Table("genres AS g").Select(selects.String(), args...).Limit(int(params.Limit)).Offset(offset)

	if params.MyFavorites != nil && *params.MyFavorites {
		claims, ok := c.Get("claims")
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "получение избранного доступно только авторизованным пользователям"})
			return
		}
		query = query.Joins("INNER JOIN user_favorite_genres AS ufg ON ufg.genre_id = g.id").
			Where("ufg.user_id = ?", claims.(*auth.Claims).ID)
	}

	if params.Query != nil {
		query = query.Where("g.name % ?", *params.Query).Order("distance ASC")
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

	var result []dto.ResponseGenreDTO

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
