package chapters

import (
	"context"
	"database/sql"
	"io"
	"log"
	"net/http"
	"slices"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/common/db/utils"
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

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 500<<20)

	var userRoles []string
	h.DB.Raw(
		`SELECT roles.name FROM roles
		INNER JOIN user_roles ON roles.id = user_roles.role_id
		WHERE user_roles.user_id = ?`, claims.ID,
	).Scan(&userRoles)

	if !slices.Contains(userRoles, "team_leader") && !slices.Contains(userRoles, "ex_team_leader") {
		c.AbortWithStatusJSON(403, gin.H{"error": "у вас недостаточно прав для добавления глав"})
		return
	}

	contentType := c.Request.Header.Get("Content-Type")
	if contentType[:19] != "multipart/form-data" {
		c.AbortWithStatusJSON(400, gin.H{"error": "тип тела запроса должен быть multipart/form-data"})
		return
	}

	reader, err := c.Request.MultipartReader()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	var name, description string

	var chapterPages ChapterPages
	chapterPages.Pages = make([][]byte, 0, 50)

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
			return
		}

		if part.FileName() != "" && part.FormName() == "pages" {
			data, err := io.ReadAll(part)
			if err != nil {
				log.Println(err)
				c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
				return
			}

			chapterPages.Pages = append(chapterPages.Pages, data)
		}

		if part.FormName() == "name" || part.FormName() == "description" {
			data, err := io.ReadAll(part)
			if err != nil {
				log.Println(err)
				c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
				return
			}

			if part.FormName() == "name" {
				name = string(data)
			} else {
				description = string(data)
			}
		}
	}

	if name == "" || len(chapterPages.Pages) == 0 {
		c.AbortWithStatusJSON(400, gin.H{"error": "в запросе не хватает страниц главы или её названия"})
		return
	}

	volumeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "id тома должен быть числом"})
		return
	}

	var titleID uint
	h.DB.Raw("SELECT title_id FROM volumes WHERE id = ?", volumeID).Scan(&titleID)
	if titleID == 0 {
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
		NumberOfPages: len(chapterPages.Pages),
		VolumeID:      sql.NullInt64{Int64: int64(volumeID), Valid: true},
		CreatorID:     claims.ID,
	}

	tx := h.DB.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	if result := tx.Create(&chapter); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	chapterPages.ChapterID = chapter.ID

	if _, err := h.ChaptersOnModerationPages.InsertOne(context.Background(), chapterPages); err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "глава успешно отправлена на модерацию"})

	conn, err := grpc.NewClient("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	client := pb.NewNotificationsClient(conn)

	if _, err := client.NotifyAboutChapterOnModeration(context.Background(), &pb.ChapterOnModeration{ID: uint64(chapter.ID), New: true}); err != nil {
		log.Println(err)
	}
}
