package chapters

import (
	"context"
	"database/sql"
	"io"
	"log"
	"slices"
	"sync"

	"github.com/Araks1255/mangacage/pkg/common/models"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ChapterPages struct {
	ChapterID uint     `bson:"chapter_id"`
	Pages     [][]byte `bson:"pages"`
}

func (h handler) CreateChapter(c *gin.Context) {
	claims := c.MustGet("claims").(*models.Claims)

	var userRoles []string
	h.DB.Raw(`SELECT roles.name FROM roles
		INNER JOIN user_roles ON roles.id = user_roles.role_id
		WHERE user_roles.user_id = ?`, claims.ID).Scan(&userRoles)

	if !slices.Contains(userRoles, "team_leader") {
		c.AbortWithStatusJSON(403, gin.H{"error": "Добавлять главы может только лидер команды"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	if len(form.Value["title"]) == 0 || len(form.Value["volume"]) == 0 || len(form.Value["name"]) == 0 || len(form.File["pages"]) == 0 {
		c.AbortWithStatusJSON(400, gin.H{"error": "в запросе не хватает названия тайтла, тома, главы или страниц главы"})
		return
	}

	title := form.Value["title"][0]
	volume := form.Value["volume"][0]
	name := form.Value["name"][0]

	var description string
	if len(form.Value["description"]) != 0 {
		description = form.Value["description"][0]
	}

	pages := form.File["pages"]

	var titleID, volumeID uint
	row := h.DB.Raw(`SELECT titles.id, volumes.id FROM titles
		INNER JOIN volumes ON titles.id = volumes.title_id
		WHERE lower(titles.name) = lower(?)
		AND lower(volumes.name) = lower(?)
		AND NOT titles.on_moderation
		AND NOT volumes.on_moderation`, title, volume).Row()

	row.Scan(&titleID, &volumeID)

	if titleID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "тайтл не найден"})
		return
	}

	if volumeID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "том не найден"})
		return
	}

	var existingChapterID uint
	h.DB.Raw("SELECT id FROM chapters WHERE volume_id = ? AND lower(name) = lower(?)", volumeID, name).Scan(&existingChapterID)
	if existingChapterID != 0 {
		c.AbortWithStatusJSON(403, gin.H{"error": "Глава уже существует"})
		return
	}

	var IsUserTeamTranslatesThisTitle bool
	h.DB.Raw(`SELECT (SELECT team_id FROM titles WHERE id = ?) = (SELECT team_id FROM users WHERE id = ?)`, titleID, claims.ID).Scan(&IsUserTeamTranslatesThisTitle)

	if !IsUserTeamTranslatesThisTitle {
		c.AbortWithStatusJSON(403, gin.H{"error": "ваша команда не переводит данный тайтл"})
		return
	}

	chapter := models.ChapterOnModeration{
		Name:          name,
		Description:   description,
		NumberOfPages: len(pages),
		VolumeID:      sql.NullInt64{Int64: int64(volumeID), Valid: true},
		CreatorID:     claims.ID,
	}

	tx := h.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	if result := tx.Create(&chapter); result.Error != nil {
		tx.Rollback()
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	chapterPages := ChapterPages{
		ChapterID: chapter.ID,
		Pages:     make([][]byte, len(pages), len(pages)),
	}

	errChan := make(chan error, len(pages))
	var wg sync.WaitGroup
	wg.Add(len(pages))

	for i := 0; i < len(pages); i++ {
		go func(index int) {
			defer wg.Done()

			page := pages[index]

			file, err := page.Open()
			if err != nil {
				errChan <- err
				file.Close()
				log.Println(err)
				return
			}
			defer file.Close()

			data, err := io.ReadAll(file)
			if err != nil {
				errChan <- err
				log.Println(err)
				return
			}

			chapterPages.Pages[index] = data

			errChan <- nil
		}(i)

	}

	wg.Wait()

	for i := 0; i < len(errChan); i++ {
		err = <-errChan
		if err != nil {
			tx.Rollback()
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}
	}

	if _, err := h.ChaptersOnModerationPages.InsertOne(context.Background(), chapterPages); err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(500, gin.H{"error": err})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "Глава успешно отправлена на модерацию"})

	conn, err := grpc.NewClient("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	client := pb.NewNotificationsClient(conn)

	if _, err := client.NotifyAboutChapterOnModeration(context.Background(), &pb.ChapterOnModeration{Name: chapter.Name, New: true}); err != nil {
		log.Println(err)
	}
}
