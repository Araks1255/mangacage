package titles

import (
	"fmt"
	"log"
	"strings"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/gin-gonic/gin"
)

type GetTitlesParams struct {
	dto.CommonParams

	PublishingStatus  *string `form:"publishingStatus"`
	TranslatingStatus *string `form:"translatingStatus"`
	Type              *string `form:"type"`

	AuthorID *uint `form:"authorId"`
	TeamID   *uint `form:"teamId"`

	YearFrom *int `form:"yearFrom"`
	YearTo   *int `form:"yearTo"`

	AgeLimitFrom *int `form:"ageLimitFrom"`
	AgeLimitTo   *int `form:"ageLimitTo"`

	ViewsFrom *int `form:"viewsFrom"`
	ViewsTo   *int `form:"viewsTo"`

	RateFrom *int `form:"rateFrom"`
	RateTo   *int `form:"rateTo"`

	ChaptersFrom *int `form:"chaptersFrom"`
	ChaptersTo   *int `form:"chaptersTo"`

	Genres []string `form:"genres"`
	Tags   []string `form:"tags"`

	FavoritedBy *uint `form:"favoritedBy" binding:"excluded_with=MyFavorites"`
	MyFavorites *bool `form:"myFavorites" binding:"excluded_with=FavoritedBy"`
	Hidden      *bool `form:"hidden"`
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

	var selects strings.Builder
	args := make([]any, 0, 3)

	selects.WriteString("t.id, t.name")

	if params.Query != nil {
		selects.WriteString(",t.name <-> ? AS name_distance, t.english_name <-> ? AS english_name_distance, t.original_name <-> ? AS original_name_distance")
		args = append(args, *params.Query, *params.Query, *params.Query)
	}

	query := h.DB.Table("titles AS t").
		Select(selects.String(), args...).
		Limit(int(params.Limit)).Offset(offset)

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

	claims, ok := c.Get("claims")

	if ok && params.MyFavorites != nil && *params.MyFavorites {
		query = query.Joins("INNER JOIN user_favorite_titles AS uft ON uft.title_id = t.id").
			Where("uft.user_id = ?", claims.(*auth.Claims).ID)
	}

	if ok && params.Hidden != nil && *params.Hidden {
		query = query.Where(
			`EXISTS(
					SELECT
						1
					FROM
						user_roles AS ur
						INNER JOIN roles AS r ON r.id = ur.role_id
					WHERE
						ur.user_id = ? AND r.name = 'admin'
				)`,
			claims.(*auth.Claims).ID,
		).
			Where("t.hidden")
	} else {
		query = query.Where("NOT t.hidden")
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

	if params.ChaptersFrom != nil {
		query = query.Where("t.number_of_chapters >= ?", params.ChaptersFrom)
	}
	if params.ChaptersTo != nil {
		query = query.Where("t.number_of_chapters <= ?", params.ChaptersTo)
	}

	if params.RateFrom != nil {
		query = query.Where(
			`CASE
				WHEN t.number_of_rates = 0 THEN 0
				ELSE t.sum_of_rates / t.number_of_rates::FLOAT
			END
				>= ?::FLOAT`,
			params.RateFrom,
		)
	}

	if params.RateTo != nil {
		query = query.Where(
			`CASE
				WHEN t.number_of_rates = 0 THEN 0
				ELSE t.sum_of_rates / t.number_of_rates::FLOAT
			END
				<= ?::FLOAT`,
			params.RateTo,
		)
	}

	if params.Query != nil {
		query = query.
			Where("t.name % ? OR t.english_name % ? OR t.original_name % ?", *params.Query, *params.Query, *params.Query).
			Order("name_distance, english_name_distance, original_name_distance ASC")
	} else {
		if params.Order != "desc" && params.Order != "asc" {
			params.Order = "asc"
		}

		switch params.Sort {
		case "views":
			query = query.Order(fmt.Sprintf("t.views %s", params.Order))

		case "createdAt":
			query = query.Order(fmt.Sprintf("t.id %s", params.Order))

		default:
			query = query.Order(fmt.Sprintf("t.name %s", params.Order))
		}
	}

	if len(params.Genres) != 0 || len(params.Tags) != 0 {
		query = query.Group("t.id")
	}

	if len(params.Genres) != 0 {
		query = query.Joins("INNER JOIN title_genres AS tg ON tg.title_id = t.id").
			Joins("INNER JOIN genres AS g ON g.id = tg.genre_id").
			Where("g.name IN ?", params.Genres).
			Having("COUNT(g.id) = ?", len(params.Genres))
	}

	if len(params.Tags) != 0 {
		query = query.Joins("INNER JOIN title_tags AS tt ON tt.title_id = t.id").
			Joins("INNER JOIN tags ON tags.id = tt.tag_id").
			Where("tags.name IN ?", params.Tags).
			Having("COUNT(tags.id) = ?", len(params.Tags))
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
