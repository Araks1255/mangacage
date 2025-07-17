package chapters

import (
	"fmt"
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
)

type GetChaptersParams struct {
	Sort  string  `form:"sort"`
	Query *string `form:"query"`
	Order string  `form:"order"`
	Page  int     `form:"page,default=1"`
	Limit uint    `form:"limit,default=20"`

	NumberOfPagesFrom *int  `form:"numberOfPagesFrom"`
	NumberOfPagesTo   *int  `form:"numberOfPagesTo"`
	ViewsFrom         *uint `form:"viewsFrom"`
	ViewsTo           *uint `form:"viewsTo"`

	Volume  *uint `form:"volume"`
	TitleID *uint `form:"titleId"`
	TeamID  *uint `form:"teamId"`

	FavoritedBy *uint `form:"favoritedBy" binding:"excluded_with=MyFavorites"`
	MyFavorites *bool `form:"myFavorites" binding:"excluded_with=FavoritedBy"`
	CreatedBy   *uint `form:"createdBy" binding:"excluded_with=MyCreations"`
	MyCreations *bool `form:"myCreations" binding:"excluded_with=CreatedBy"`

	Viewed *bool `form:"viewed"`
}

func (h handler) GetChapters(c *gin.Context) {
	var params GetChaptersParams

	if err := c.ShouldBindQuery(&params); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	query := h.DB.Table("chapters AS c").
		Select("c.*, teams.name AS team, t.name AS title").
		Joins("INNER JOIN titles AS t ON t.id = c.title_id").
		Joins("LEFT JOIN teams ON c.team_id = teams.id")

	if params.TitleID != nil {
		query = query.Where("t.id = ?", params.TitleID)
	}

	if params.TeamID != nil {
		query = query.Where("c.team_id = ?", params.TeamID)
	}

	if params.Volume != nil {
		query = query.Where("c.volume = ?", params.Volume)
	}

	if params.Query != nil {
		query = query.Where("lower(c.name) ILIKE lower(?)", fmt.Sprintf("%%%s%%", *params.Query))
	}

	if params.ViewsFrom != nil {
		query = query.Where("c.views >= ?", params.ViewsFrom)
	}
	if params.ViewsTo != nil {
		query = query.Where("c.views <= ?", params.ViewsTo)
	}

	if params.NumberOfPagesFrom != nil {
		query = query.Where("c.number_of_pages >= ?", params.NumberOfPagesFrom)
	}
	if params.NumberOfPagesTo != nil {
		query = query.Where("c.number_of_pages <= ?", params.NumberOfPagesTo)
	}

	if params.FavoritedBy != nil {
		query = query.Joins("INNER JOIN user_favorite_chapters AS ufc ON ufc.chapter_id = c.id").
			Joins("INNER JOIN users AS u ON u.id = ufc.user_id").
			Where("u.id = ?", *params.FavoritedBy).
			Where("u.visible")
	}

	if params.CreatedBy != nil {
		query = query.Joins("INNER JOIN users AS creators ON creators.id = c.creator_id").
			Where("creators.id = ?", *params.CreatedBy).
			Where("creators.visible")
	}

	if params.MyCreations != nil && *params.MyCreations {
		claims, ok := c.Get("claims")
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "получение своих публикаций доступно только авторизованным пользователям"})
			return
		}
		query = query.Where("c.creator_id = ?", claims.(*auth.Claims).ID)
	}

	if params.MyFavorites != nil && *params.MyFavorites {
		claims, ok := c.Get("claims")
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "получение избранного доступно только авторизованным пользователям"})
			return
		}
		query = query.Joins("INNER JOIN user_favorite_chapters AS ufc ON ufc.chapter_id = c.id").
			Where("ufc.user_id = ?", claims.(*auth.Claims).ID)
	}

	if params.Viewed != nil && *params.Viewed {
		claims, ok := c.Get("claims")
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "получение истории чтения глав доступно только авторизованным пользователям"})
			return
		}
		query = query.Joins("INNER JOIN user_viewed_chapters AS uvc ON uvc.chapter_id = c.id").
			Where("uvc.user_id = ?", claims.(*auth.Claims).ID)
	}

	if params.Order != "desc" && params.Order != "asc" {
		params.Order = "desc"
	}

	switch params.Sort {
	case "views":
		query = query.Order(fmt.Sprintf("c.views %s", params.Order))

	case "numberOfPages":
		query = query.Order(fmt.Sprintf("c.number_of_pages %s", params.Order))

	case "createdAt":
		query = query.Order(fmt.Sprintf("c.id %s", params.Order))

	case "readedAt":
		if params.Viewed != nil && *params.Viewed {
			query = query.Order(fmt.Sprintf("uvc.created_at %s", params.Order))
			break
		}
	default:
		query = query.Order(fmt.Sprintf("c.name %s", params.Order))
	}

	offset := (params.Page - 1) * int(params.Limit)
	if offset < 0 {
		offset = 0
	}

	query = query.Limit(int(params.Limit)).Offset(offset)

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
