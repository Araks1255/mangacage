package titles

import (
	"context"
	"fmt"
	"io"
	"log"
	"strings"

	"github.com/Araks1255/mangacage/pkg/common/models"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"gorm.io/gorm"
)

func (h handler) CreateTitle(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	name := strings.ToLower(form.Value["name"][0])
	description := strings.ToLower(form.Value["description"][0])
	author := strings.ToLower(form.Value["author"][0])

	genres := form.Value["genres"]

	cover, err := c.FormFile("cover")
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	var existingTitleID uint
	h.DB.Raw("SELECT id FROM titles WHERE name = ?", name).Scan(&existingTitleID)
	if existingTitleID != 0 {
		c.AbortWithStatusJSON(403, gin.H{"error": "Тайтл уже существует"})
		return
	}

	var authorID uint
	h.DB.Raw("SELECT id FROM authors WHERE name = ?", author).Scan(&authorID)
	if authorID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "Автор не найден"})
		return
	}

	title := models.Title{
		Name:         name,
		Description:  description,
		AuthorID:     authorID,
		CreatorID:    claims.ID,
		OnModeration: true,
	}

	tx := h.DB.Begin()

	if result := tx.Create(&title); result.Error != nil {
		tx.Rollback()
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if err := AddGenresToTitle(title.ID, genres, tx); err != nil {
		tx.Rollback()
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	var titleCover struct {
		TitleID uint   `bson:"title_id"`
		Cover   []byte `bson:"cover"`
	}

	file, err := cover.Open()
	if err != nil {
		tx.Rollback()
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	defer file.Close()

	titleCover.Cover, err = io.ReadAll(file)
	if err != nil {
		tx.Rollback()
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if _, err := h.Collection.InsertOne(context.Background(), titleCover); err != nil {
		tx.Rollback()
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

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

func AddGenresToTitle(titleID uint, genres []string, tx *gorm.DB) error {
	query := "INSERT INTO title_genres (title_id, genre_id) VALUES"

	for i := 0; i < len(genres); i++ {
		query += fmt.Sprintf(" (%d, (SELECT id FROM genres WHERE name = '%s')),", titleID, strings.ToLower(genres[i]))
	}

	query = strings.TrimSuffix(query, ",")

	if result := tx.Exec(query); result.Error != nil {
		return result.Error
	}

	return nil
}
