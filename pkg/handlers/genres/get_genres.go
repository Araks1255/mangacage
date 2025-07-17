package genres

import (
	"fmt"
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
)

type getGenresParams struct {
	Query *string `form:"query"`
	Order string  `form:"order"`
	Sort  string  `form:"sort"`
	Page  int     `form:"page,default=1"`
	Limit uint    `form:"limit,default=20"`

	FavoritedBy *uint `form:"favoritedBy" binding:"excluded_with=MyFavorites"`
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

	query := h.DB.Table("genres AS g").Select("g.*").Limit(int(params.Limit)).Offset(offset)

	if params.Query != nil {
		query = query.Where("lower(g.name) ILIKE lower(?)", fmt.Sprintf("%%%s%%", *params.Query))
	}

	if params.FavoritedBy != nil {
		query = query.Joins("INNER JOIN user_favorite_genres AS ufg ON ufg.genre_id = g.id").
			Joins("INNER JOIN users AS u ON ufg.user_id = u.id").
			Where("u.id = ?", params.FavoritedBy).
			Where("u.visible")
	}

	if params.MyFavorites != nil && *params.MyFavorites {
		claims, ok := c.Get("claims")
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "получение избранного доступно только авторизованным пользователям"})
			return
		}
		query = query.Joins("INNER JOIN user_favorite_genres AS ufg ON ufg.genre_id = g.id").
			Where("ufg.user_id = ?", claims.(*auth.Claims).ID)
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
