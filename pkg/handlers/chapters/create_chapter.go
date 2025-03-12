package chapters

import (
	"context"
	"io"
	"log"
	"slices"

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

	title := c.Param("title")
	volume := c.Param("volume")

	name := form.Value["name"][0]
	description := form.Value["description"][0]

	pages := form.File["pages"]
	if len(pages) == 0 {
		c.AbortWithStatusJSON(400, gin.H{"error": "отсутствуют страницы главы"})
		return
	}

	var titleID uint
	h.DB.Raw("SELECT id FROM titles WHERE lower(name) = lower(?)", title).Scan(&titleID)
	if titleID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "Тайтл не найден"})
		return
	}

	var volumeID uint
	h.DB.Raw("SELECT id FROM volumes WHERE lower(name) = lower(?) AND title_id = ?", volume, titleID).Scan(&volumeID)
	if volumeID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "Том не найден"})
		return
	}

	var IsUserTeamTranslatesThisTitle bool
	h.DB.Raw(`SELECT CAST(
		CASE WHEN ? = 
		(SELECT titles.id FROM titles
		INNER JOIN teams ON titles.team_id = teams.id
		INNER JOIN users ON teams.id = users.team_id
		WHERE users.id = ?)
		THEN true ELSE false END AS BOOLEAN)`, titleID, claims.ID).Scan(&IsUserTeamTranslatesThisTitle)

	if !IsUserTeamTranslatesThisTitle {
		c.AbortWithStatusJSON(403, gin.H{"error": "Ваша команда не переводит данный тайтл"})
		return
	}

	var existingChapterID uint
	h.DB.Raw("SELECT id FROM chapters WHERE volume_id = ? AND lower(name) = lower(?)", volumeID, name).Scan(&existingChapterID)
	if existingChapterID != 0 {
		c.AbortWithStatusJSON(403, gin.H{"error": "Глава уже существует"})
		return
	}

	chapter := models.Chapter{
		Name:          name,
		Description:   description,
		NumberOfPages: len(pages),
		VolumeID:      volumeID,
		OnModeration:  true,
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

	for i := 0; i < len(pages); i++ {
		page := pages[i]

		file, err := page.Open()
		if err != nil {
			tx.Rollback()
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}

		data, err := io.ReadAll(file)
		if err != nil {
			tx.Rollback()
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}

		chapterPages.Pages[i] = data // Это не совсем безопасно, но почему-то при использовании append слайс расширялся до 6 элементов, при этом первые 3 оставались незаполненными, а это очень плохо
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

	if _, err := client.NotifyAboutNewChapterOnModeration(context.Background(), &pb.ChapterOnModeration{TitleName: title, ChapterName: name}); err != nil {
		log.Println(err)
	}
}
