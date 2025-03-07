package search

import (
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func (h handler) Search(c *gin.Context) {
	table := strings.ToLower(c.Param("type"))
	query := strings.ToLower(c.Param("query"))

	switch table {
	case "titles":
		type result struct {
			Title  string `gorm:"column:name"`
			Author string `gorm:"column:name"`
		}

		var results []result
		h.DB.Raw(`SELECT titles.name, authors.name FROM titles
			INNER JOIN authors ON titles.author_id = authors.id
			WHERE titles.name ILIKE ? AND NOT titles.on_moderation`,
			fmt.Sprintf("%%%s%%", query),
		).Scan(&results)

		if len(results) == 0 {
			c.AbortWithStatusJSON(404, gin.H{"error": "не найдено тайтлов по вашему запросу"})
			return
		}

		response := make(map[int]result, len(results))
		for i := 0; i < len(results); i++ {
			response[i] = results[i]
		}

		c.JSON(200, &response)

	case "volumes":
		type result struct {
			Title  string `gorm:"column:name"`
			Volume string `gorm:"column:name"`
		}

		var results []result
		h.DB.Raw(`SELECT titles.name, volumes.name FROM volumes
		INNER JOIN titles ON volumes.title_id = titles.id
		WHERE volumes.name ILIKE ? AND NOT volumes.on_moderation`,
			fmt.Sprintf("%%%s%%", query),
		).Scan(&results)

		if len(results) == 0 {
			c.AbortWithStatusJSON(404, gin.H{"error": "не найдено томов по вашему запросу"})
			return
		}

		response := make(map[int]result, len(results))
		for i := 0; i < len(results); i++ {
			response[i] = results[i]
		}

		c.JSON(200, &response)

	case "chapters":
		type result struct {
			Title   string `gorm:"column:name"`
			Volume  string `gorm:"column:name"`
			Chapter string `gorm:"column:name"`
		}

		var results []result
		h.DB.Raw(`SELECT titles.name, volumes.name, chapters.name FROM chapters
			INNER JOIN volumes ON chapters.volume_id = volumes.id
			INNER JOIN titles ON volumes.title_id = titles.id
			WHERE chapters.name ILIKE ? AND NOT chapters.on_moderation`,
			fmt.Sprintf("%%%s%%", query),
		).Scan(&results)

		if len(results) == 0 {
			c.AbortWithStatusJSON(404, gin.H{"error": "не найдено глав по вашему запросу"})
			return
		}

		response := make(map[int]result, len(results))
		for i := 0; i < len(results); i++ {
			response[i] = results[i]
		}

		c.JSON(200, &response)

	case "teams":
		var teams []string
		h.DB.Raw("SELECT name FROM teams WHERE name ILIKE ?", fmt.Sprintf("%%%s%%", query)).Scan(&teams)
		if len(teams) == 0 {
			c.AbortWithStatusJSON(404, gin.H{"error": "не найдено команд перевода по вашему запросу"})
			return
		}

		response := make(map[int]string, len(teams))
		for i := 0; i < len(teams); i++ {
			response[i] = teams[i]
		}

		c.JSON(200, &response)

	case "authors":
		var authors []string
		h.DB.Raw("SELECT name FROM authors WHERE name ILIKE ?", fmt.Sprintf("%%%s%%", query)).Scan(&authors)
		if len(authors) == 0 {
			c.AbortWithStatusJSON(404, gin.H{"error": "не найдено авторов по вашему запросу"})
			return
		}

		response := make(map[int]string, len(authors))
		for i := 0; i < len(authors); i++ {
			response[i] = authors[i]
		}

		c.JSON(200, &response)

	default:
		c.AbortWithStatusJSON(400, gin.H{"error": "некорректная область поиска"})
		return
	}
}
