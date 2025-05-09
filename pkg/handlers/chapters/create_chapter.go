package chapters

import (
	"database/sql"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
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

	if err := h.DB.Raw(
		`SELECT r.name FROM roles AS r
		INNER JOIN user_roles AS ur ON ur.role_id = r.id
		WHERE ur.user_id = ?`, claims.ID,
	).Scan(&userRoles).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if !slices.Contains(userRoles, "team_leader") && !slices.Contains(userRoles, "ex_team_leader") {
		c.AbortWithStatusJSON(403, gin.H{"error": "у вас недостаточно прав для добавления глав"})
		return
	}

	volumeID, err := strconv.ParseUint(c.Param("id"), 10, 64)
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

	if len(pages) == 0 {
		c.AbortWithStatusJSON(400, gin.H{"error": "в запросе не хватает страниц главы"})
		return
	}
	if name == "" {
		c.AbortWithStatusJSON(400, gin.H{"error": "в запросе не хватает названия главы"})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var check struct {
		VolumeID      sql.NullInt64
		ChapterExists bool
	}

	if err = tx.Raw(
		`SELECT
			(SELECT id FROM volumes WHERE id = ?) AS volume_id,
			EXISTS(SELECT 1 FROM chapters WHERE lower(name) = lower(?)) AS chapter_exists`,
		volumeID, name,
	).Scan(&check).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if !check.VolumeID.Valid {
		c.AbortWithStatusJSON(404, gin.H{"error": "том не найден"})
		return
	}
	if check.ChapterExists {
		c.AbortWithStatusJSON(409, gin.H{"error": "глава с таким названием уже существует в этом томе"})
		return
	}

	newChapter := models.ChapterOnModeration{
		Name:          sql.NullString{String: name, Valid: true},
		Description:   description,
		NumberOfPages: len(pages),
		VolumeID:      uint(check.VolumeID.Int64),
		CreatorID:     claims.ID,
	}

	err = tx.Create(&newChapter).Error

	if err != nil {
		if dbErrors.IsUniqueViolation(err, constraints.UniChapterVolume) {
			c.AbortWithStatusJSON(409, gin.H{"error": "глава с таким названием уже ожидает модерации в этом томе"})
		} else {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}
		return
	}

	chapterPages := ChapterPages{
		ChapterOnModerationID: newChapter.ID,
		Pages:                 pages,
	}

	if _, err := h.ChaptersOnModerationPages.InsertOne(c.Request.Context(), chapterPages); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "глава успешно отправлена на модерацию"})

	if _, err := h.NotificationsClient.NotifyAboutChapterOnModeration(c.Request.Context(), &pb.ChapterOnModeration{ID: uint64(newChapter.ID), New: true}); err != nil {
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
