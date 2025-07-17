package titles

import (
	"fmt"
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type GetTitlesParams struct {
	Sort  string  `form:"sort"`
	Order string  `form:"order"`
	Query *string `form:"query"`
	Page  int     `form:"page,default=1"`
	Limit uint    `form:"limit,default=20"`

	PublishingStatus  *string `form:"publishingStatus"`
	TranslatingStatus *string `form:"translatingStatus"`
	Type              *string `form:"type"`

	AuthorID *uint `form:"authorId"`
	TeamID   *uint `form:"teamId"`

	YearFrom     *int `form:"yearFrom"`
	YearTo       *int `form:"yearTo"`
	AgeLimitFrom *int `form:"ageLimitFrom"`
	AgeLimitTo   *int `form:"ageLimitTo"`
	ViewsFrom    *int `form:"viewsFrom"`
	ViewsTo      *int `form:"viewsTo"`
	RateFrom     *int `form:"rateFrom"`
	RateTo       *int `form:"rateTo"`
	ChaptersFrom *int `form:"chaptersFrom"`
	ChaptersTo   *int `form:"chaptersTo"`

	Genres []string `form:"genres"`
	Tags   []string `form:"tags"`

	FavoritedBy *uint `form:"favoritedBy" binding:"excluded_with=MyFavorites"`
	MyFavorites *bool `form:"myFavorites" binding:"excluded_with=FavoritedBy"`
}

func (h handler) GetTitles(c *gin.Context) {
	var params GetTitlesParams

	if err := c.ShouldBindQuery(&params); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	offset := (params.Page - 1) * int(params.Limit)
	if offset < 0 {
		offset = 0
	}

	query := h.DB.Table("titles AS t").
		Select("t.*, a.name AS author").
		Joins("INNER JOIN authors AS a ON a.id = t.author_id").
		Limit(int(params.Limit)).Offset(offset)

	if params.Query != nil {
		query = query.Where(
			"lower(t.name) ILIKE lower(?) OR lower(t.english_name) ILIKE lower(?) OR t.original_name ILIKE ?",
			fmt.Sprintf("%%%s%%", *params.Query), fmt.Sprintf("%%%s%%", *params.Query), fmt.Sprintf("%%%s%%", *params.Query),
		)
	}

	if params.PublishingStatus != nil {
		query = query.Where("t.publishing_status = ?", params.PublishingStatus)
	}

	if params.TranslatingStatus != nil {
		query = query.Where("t.translating_status = ?", params.TranslatingStatus)
	}

	if params.Type != nil {
		query = query.Where("t.type = ?", params.Type)
	}

	if params.AuthorID != nil {
		query = query.Where("t.author_id = ?", params.AuthorID)
	}

	if params.TeamID != nil {
		query = query.Joins("INNER JOIN title_teams ON title_teams.title_id = t.id").Where("title_teams.team_id = ?", params.TeamID)
	}

	if params.FavoritedBy != nil {
		query = query.Joins("INNER JOIN user_favorite_titles AS uft ON uft.title_id = t.id").
			Joins("INNER JOIN users AS u ON uft.user_id = u.id").
			Where("user_id = ?", params.FavoritedBy).
			Where("u.visible")
	}

	if params.MyFavorites != nil && *params.MyFavorites {
		claims, ok := c.Get("claims")
		if !ok {
			c.AbortWithStatusJSON(401, gin.H{"error": "получение избранного доступно только авторизованным пользователям"})
			return
		}
		query = query.Joins("INNER JOIN user_favorite_titles AS uft ON uft.title_id = t.id").
			Where("uft.user_id = ?", claims.(*auth.Claims).ID)
	}

	if params.YearFrom != nil {
		query = query.Where("t.year_of_release >= ?", params.YearFrom)
	}
	if params.YearTo != nil {
		query = query.Where("t.year_of_release <= ?", params.YearTo)
	}

	if params.AgeLimitFrom != nil {
		query = query.Where("t.age_limit >= ?", params.AgeLimitFrom)
	}
	if params.AgeLimitTo != nil {
		query = query.Where("t.age_limit <= ?", params.AgeLimitTo)
	}

	if params.ViewsFrom != nil {
		query = query.Where("t.views >= ?", params.ViewsFrom)
	}
	if params.ViewsTo != nil {
		query = query.Where("t.views <= ?", params.ViewsTo)
	}

	if params.RateFrom != nil {
		query = query.Where("t.sum_of_rates/t.number_of_rates::FLOAT >= ?::FLOAT", params.RateFrom)
	}
	if params.RateTo != nil {
		query = query.Where("t.sum_of_rates/t.number_of_rates::FLOAT <= ?::FLOAT", params.RateTo)
	}

	if params.Order != "desc" && params.Order != "asc" {
		params.Order = "desc"
	}

	switch params.Sort {
	case "views":
		query = query.Order(fmt.Sprintf("t.views %s", params.Order))

	case "createdAt":
		query = query.Order(fmt.Sprintf("t.id %s", params.Order))

	default:
		query = query.Order("t.name DESC")
	}

	if params.Genres != nil {
		query = query.Where(
			`(SELECT ARRAY(
				SELECT g.name FROM genres AS g
				INNER JOIN title_genres AS tg ON tg.genre_id = g.id
				WHERE tg.title_id = t.id
			))::TEXT[] @> ?::TEXT[]`,
			pq.Array(params.Genres),
		)
	}

	if params.Tags != nil {
		query = query.Where(
			`(SELECT ARRAY(
				SELECT tags.name FROM tags
				INNER JOIN title_tags AS tt ON tt.tag_id = tags.id
				WHERE tt.title_id = t.id
			))::TEXT[] @> ?::TEXT[]`,
			pq.Array(params.Tags),
		)
	}

	var result []dto.ResponseTitleDTO

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
