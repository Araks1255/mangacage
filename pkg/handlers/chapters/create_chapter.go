package chapters

import (
	"context"
	"database/sql"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
)

type ChapterPages struct {
	ChapterOnModerationID uint     `bson:"chapter_on_moderation_id"`
	Pages                 [][]byte `bson:"pages"`
}

func (h handler) CreateChapter(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

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

	desiredVolumeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id тома"})
		return
	}

	c.Request.Body = http.MaxBytesReader(c.Writer, c.Request.Body, 500<<20)

	contentType := c.Request.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "multipart/form-data") {
		c.AbortWithStatusJSON(400, gin.H{"error": "тело запроса должно иметь тип multipart/form-data"})
		return
	}

	reader, err := c.Request.MultipartReader()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	name, description, pages, err := parseCreateChapterBody(reader)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if name == "" || len(pages) == 0 {
		c.AbortWithStatusJSON(400, gin.H{"error": "в запросе не хватает названия главы или её страниц"})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var existing struct {
		VolumeID              uint
		ChapterOnModerationID uint
		ChapterID             uint
	}

	tx.Raw(
		`SELECT
			v.id AS volume_id,
			com.id AS chapter_on_moderation_id,
			c.id AS chapter_id
		FROM
			volumes AS v
			LEFT JOIN chapters_on_moderation AS com ON v.id = com.volume_id AND lower(com.name) = lower(?)
			LEFT JOIN chapters AS c ON v.id = c.volume_id AND lower(c.name) = lower(?)
		WHERE
			v.id = ?
		LIMIT 1`,
		name, name, desiredVolumeID,
	).Scan(&existing)

	if existing.VolumeID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "том не найден"})
		return
	}
	if existing.ChapterOnModerationID != 0 {
		c.AbortWithStatusJSON(409, gin.H{"error": "глава с таким названием уже ожидает модерации в этом томе"})
		return
	}
	if existing.ChapterID != 0 {
		c.AbortWithStatusJSON(409, gin.H{"error": "глава с таким названием уже существует в этом томе"})
		return
	}

	newChapter := models.ChapterOnModeration{
		Name:          name,
		Description:   description,
		NumberOfPages: len(pages),
		VolumeID:      sql.NullInt64{Int64: int64(existing.VolumeID), Valid: true},
		CreatorID:     claims.ID,
	}

	if result := tx.Create(&newChapter); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	chapterPages := ChapterPages{
		ChapterOnModerationID: newChapter.ID,
		Pages:                 pages,
	}

	if _, err := h.ChaptersOnModerationPages.InsertOne(context.Background(), chapterPages); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "глава успешно отправлена на модерацию"})

	if _, err := h.NotificationsClient.NotifyAboutChapterOnModeration(context.Background(), &pb.ChapterOnModeration{ID: uint64(newChapter.ID), New: true}); err != nil {
		log.Println(err)
	}
}

func parseCreateChapterBody(r *multipart.Reader) (name, description string, pages [][]byte, parsingError error) { // Тут я сделал ручной парсинг формы, потому-что при использовании *gin.Context.MultipartForm() все файлы из тела запроса сохраняются в виде *multipart.FileHeader (содержащих сами файлы), а потом, при их чтении, полученные срезы байт вновь сохраняются в оперативной памяти, и нагрузка на неё идёт двойная (хотя тут вообще по-хорошему надо сделать постепенную вставку в mongoDB, чтобы страницы в целом в памяти не находились все разом. Но этим я займусь позже)
	pages = make([][]byte, 0, 45)
	for {
		part, err := r.NextPart()

		if err == io.EOF {
			break
		}
		if err != nil {
			return "", "", nil, err
		}

		if part.FormName() == "name" {
			data, err := io.ReadAll(part)
			if err != nil {
				return "", "", nil, err
			}
			name = string(data)
		}

		if part.FormName() == "description" {
			data, err := io.ReadAll(part)
			if err != nil {
				return "", "", nil, err
			}
			description = string(data)
		}

		if part.FileName() != "" && part.FormName() == "pages" {
			page, err := io.ReadAll(part)
			if err != nil {
				return "", "", nil, err
			}
			pages = append(pages, page)
		}
	}

	return name, description, pages, nil
}
