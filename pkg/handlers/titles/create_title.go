package titles

import (
	"context"
	"fmt"
	"log"
	"slices"
	"strings"

	"github.com/Araks1255/mangacage/pkg/common/models"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
)

func (h handler) CreateTitle(c *gin.Context) {
	claims, ok := c.MustGet("claims").(*models.Claims)
	if !ok {
		log.Println("Не удалось привести клаймсы к типу")
		c.AbortWithStatusJSON(500, gin.H{"error": "Внутренняя ошибка, извините"})
		return
	}

	var requestBody struct {
		Name        string   `json:"name" binding:"required"`
		Description string   `json:"description"`
		Author      string   `json:"author" binding:"required"`
		Genres      []string `json:"genres" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	transaction := h.DB.Begin()

	var authorID uint
	transaction.Raw("SELECT id FROM authors WHERE name = ?", requestBody.Author).Scan(&authorID)
	if authorID == 0 {
		transaction.Rollback()
		c.AbortWithStatusJSON(404, gin.H{"error": "Автор не найден"})
		return
	}

	title := models.Title{
		Name:        requestBody.Name,
		Description: requestBody.Description,
		AuthorID:    authorID,
		CreatorID:   claims.ID,
	}

	var userRoles []string
	transaction.Raw("SELECT roles.name FROM roles "+
		"INNER JOIN user_roles ON roles.id = user_roles.role_id "+
		"INNER JOIN users ON user_roles.user_id = users.id "+ // Я знаю, что конкатенация в го процесс сложный, но не хочу делать длиннющую строку
		"WHERE users.id = ?", claims.ID).Scan(&userRoles)

	if IsUserAdmin := slices.Contains(userRoles, "admin"); IsUserAdmin == true {
		title.OnModeration = false
	} else {
		title.OnModeration = true
	}

	if result := transaction.Create(&title); result.Error != nil {
		transaction.Rollback()
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": "Не удалось создать тайтл"})
		return
	}

	if err := AddGenresToTitle(title.ID, requestBody.Genres, transaction); err != nil {
		transaction.Rollback()
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": "Не удалось добавить жанры к тайтлу"})
		return
	}

	transaction.Commit()

	c.JSON(201, gin.H{"success": "Тайтл успешно создан"})

	conn, err := grpc.NewClient("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	client := pb.NewNotificationsClient(conn)

	if _, err = client.NotifyAboutNewTitleOnModeration(context.Background(), &pb.TitleOnModeration{TitleName: title.Name}); err != nil {
		log.Println(err)
	}
}

func AddGenresToTitle(titleID uint, genres []string, transaction *gorm.DB) error {
	query := "INSERT INTO title_genres (title_id, genre_id) VALUES"

	for i := 0; i < len(genres); i++ {
		query += fmt.Sprintf(" (%d, (SELECT id FROM genres WHERE name = '%s')),", titleID, genres[i])
	}

	query = strings.TrimSuffix(query, ",")

	if result := transaction.Exec(query); result.Error != nil {
		return result.Error
	}

	return nil
}
