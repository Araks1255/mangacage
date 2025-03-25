package titles

import (
	"context"
	"database/sql"
	"io"
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (h handler) CreateTitle(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	if len(form.Value["name"]) == 0 || len(form.Value["author"]) == 0 || len(form.Value["genres"]) == 0 {
		c.AbortWithStatusJSON(400, gin.H{"error": "в запросе недостаточно данных"})
		return
	}

	name := form.Value["name"][0]
	author := form.Value["author"][0]

	var description string
	if len(form.Value["description"]) != 0 {
		description = form.Value["description"][0]
	}

	genres := pq.StringArray(form.Value["genres"])

	cover, err := c.FormFile("cover")
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	var titleCover struct {
		TitleOnModerationID uint   `bson:"title_on_moderation_id"`
		Cover               []byte `bson:"cover"`
	}

	errChan := make(chan error)

	go func() {
		file, err := cover.Open()
		if err != nil {
			log.Println(err)
			errChan <- err
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			log.Println(err)
			errChan <- err
			return
		}

		titleCover.Cover = data

		errChan <- nil
	}()

	var existingTitleID uint
	h.DB.Raw("SELECT id FROM titles WHERE lower(name) = lower(?)", name).Scan(&existingTitleID)
	if existingTitleID != 0 {
		c.AbortWithStatusJSON(403, gin.H{"error": "тайтл уже существует"})
		return
	}

	var existingTitleOnModerationID uint
	h.DB.Raw("SELECT id FROM titles_on_moderation WHERE lower(name) = lower(?)", name).Scan(&existingTitleOnModerationID)
	if existingTitleOnModerationID != 0 {
		c.AbortWithStatusJSON(403, gin.H{"error": "тайтл с таким названием уже находится на модерации"})
		return
	}

	var authorID uint
	h.DB.Raw("SELECT id FROM authors WHERE lower(name) = lower(?)", author).Scan(&authorID)
	if authorID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "автор не найден"})
		return
	}

	title := models.TitleOnModeration{
		Name:        name,
		Description: description,
		CreatorID:   claims.ID,
		AuthorID:    sql.NullInt64{Int64: int64(authorID), Valid: true},
		Genres:      genres,
	}

	tx := h.DB.Begin()
	if r := recover(); r != nil {
		tx.Rollback()
		panic(r)
	}

	if result := tx.Create(&title); result.Error != nil {
		tx.Rollback()
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if err = <-errChan; err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
	}

	titleCover.TitleOnModerationID = title.ID

	if _, err := h.Collection.InsertOne(context.Background(), titleCover); err != nil {
		tx.Rollback()
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "тайтл успешно отправлен на модерацию"})

	conn, err := grpc.NewClient("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	client := pb.NewNotificationsClient(conn)

	if _, err = client.NotifyAboutTitleOnModeration(context.Background(), &pb.TitleOnModeration{Name: title.Name, New: true}); err != nil {
		log.Println(err)
	}
}
