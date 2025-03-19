package chapters

import (
	"context"
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
	h.DB.Raw("SELECT roles.name FROM roles "+
		"INNER JOIN user_roles ON roles.id = user_roles.role_id "+
		"INNER JOIN users ON user_roles.user_id = users.id "+
		"WHERE users.id = ?", claims.ID).Scan(&userRoles)

	if IsUserTeamOwner := slices.Contains(userRoles, "team_leader"); !IsUserTeamOwner {
		c.AbortWithStatusJSON(403, gin.H{"error": "Добавлять главы может только лидер команды"})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	if len(form.Value["name"]) == 0 {
		c.AbortWithStatusJSON(400, gin.H{"error": "в запросе отсутствует название главы"})
		return
	}

	title := c.Param("title")
	volume := c.Param("volume")

	name := form.Value["name"][0]

	var description string
	if len(form.Value["description"]) != 0 {
		description = form.Value["description"][0]
	}

	pages := form.File["pages"]
	if len(pages) == 0 {
		c.AbortWithStatusJSON(400, gin.H{"error": "отсутствуют страницы главы"})
		return
	}

	var titleID, volumeID uint
	row := h.DB.Raw(`SELECT titles.id, volumes.id FROM titles
		INNER JOIN volumes ON titles.id = volumes.title_id
		WHERE lower(titles.name) = lower(?)
		AND lower(volumes.name) = lower(?)
		AND NOT titles.on_moderation
		AND NOT volumes.on_moderation`, title, volume).Row()

	if err = row.Scan(&titleID, &volumeID); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

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
	h.DB.Raw(`SELECT ? = ANY(ARRAY(SELECT titles.id FROM titles
		INNER JOIN teams ON titles.team_id = teams.id
		INNER JOIN users ON teams.id = users.team_id
		WHERE users.id = ?))`, titleID, claims.ID).
		Scan(&IsUserTeamTranslatesThisTitle)

	if !IsUserTeamTranslatesThisTitle {
		c.AbortWithStatusJSON(403, gin.H{"error": "Ваша команда не переводит данный тайтл"})
		return
	}

	chapter := models.Chapter{
		Name:          name,
		Description:   description,
		NumberOfPages: len(pages),
		VolumeID:      volumeID,
	}

	tx := h.DB.Begin()

	if result := tx.Create(&chapter); result.Error != nil { // Заменить это на транзакцию, а коммитить только после выгрузки страниц
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

	if _, err := h.Collection.InsertOne(context.Background(), chapterPages); err != nil {
		tx.Rollback()
		c.AbortWithStatusJSON(500, gin.H{"error": err})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "Глава успешно создана"})

	conn, err := grpc.NewClient("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	client := pb.NewNotificationsClient(conn)

	if _, err := client.NotifyAboutChapterOnModeration(context.Background(), &pb.ChapterOnModeration{TitleName: title, Name: name}); err != nil {
		log.Println(err)
	}
}
