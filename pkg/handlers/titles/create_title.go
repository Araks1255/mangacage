package titles

import (
	"fmt"
	"log"
	"strings"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h handler) CreateTitle(c *gin.Context) {
	// claims, ok := c.MustGet("claims").(*models.Claims)
	// if !ok {
	// 	log.Println("Не удалось привести клаймсы к типу")
	// 	c.AbortWithStatusJSON(500, gin.H{"error": "Внутренняя ошибка, извините"})
	// 	return
	// }

	var requestBody struct {
		Name        string   `json:"name" binding:"required"`
		Description string   `json:"description"`
		Genres      []string `json:"genres" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	title := models.Title{
		Name:        requestBody.Name,
		Description: requestBody.Description,
	}

	if result := h.DB.Create(&title); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": "Не удалось создать тайтл"})
		return
	}

	// genresIDs := make([]uint, len(requestBody.Genres), len(requestBody.Genres))
	// var genreID uint
	// for i := 0; i < len(requestBody.Genres); i++ {
	// 	h.DB.Raw("SELECT id FROM genres WHERE name = ?", requestBody.Genres[i]).Scan(&genreID)
	// 	genresIDs[i] = genreID
	// }

	if err := AddGenresToTitle(title.ID, requestBody.Genres, h.DB); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": "Не удалось добавить жанры к тайтлу"})
		return
	}

	c.JSON(201, gin.H{"success": "Тайтл успешно создан"})
}

func AddGenresToTitle(titleID uint, genres []string, db *gorm.DB) error {
	query := "INSERT INTO title_genres (title_id, genre_id) VALUES"
	for i := 0; i < len(genres); i++ {
		query += fmt.Sprintf(" (%d, (SELECT id FROM genres WHERE name = '%s')),", titleID, genres[i])
	}
	query = strings.TrimSuffix(query, ",")

	if result := db.Exec(query); result.Error != nil {
		return result.Error
	}

	return nil
}
