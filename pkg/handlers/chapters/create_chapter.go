package chapters

import (
	"database/sql"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/common/models/mongo"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
)

func (h handler) CreateChapter(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

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
		UserTeamID       *uint
		DoesChapterExist bool
	}

	if err = tx.Raw(
		`SELECT
			(
				SELECT team_id FROM title_teams WHERE title_id = (
					SELECT title_id FROM volumes WHERE id = ?
				) AND team_id = (
					SELECT team_id FROM users WHERE id = ?
				)
			) AS user_team_id,
			EXISTS(SELECT 1 FROM chapters WHERE lower(name) = lower(?) AND volume_id = ?) AS does_chapter_exist`,
		volumeID, claims.ID, name, volumeID,
	).Scan(&check).Error; err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if check.UserTeamID == nil {
		c.AbortWithStatusJSON(404, gin.H{"error": "том не найден среди томов тайтлов, переводимых вашей командой"})
		return
	}
	if check.DoesChapterExist {
		c.AbortWithStatusJSON(409, gin.H{"error": "глава с таким названием уже существует в этом томе"})
		return
	}

	newChapter := models.ChapterOnModeration{
		Name:          sql.NullString{String: name, Valid: true},
		Description:   description,
		NumberOfPages: len(pages),
		VolumeID:      uint(volumeID),
		CreatorID:     claims.ID,
		TeamID:        *check.UserTeamID,
	}

	err = tx.Create(&newChapter).Error

	if err != nil {
		if dbErrors.IsUniqueViolation(err, constraints.UniqChapterVolume) {
			c.AbortWithStatusJSON(409, gin.H{"error": "глава с таким названием уже ожидает модерации в этом томе"})
		} else {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}
		return
	}

	chapterPages := mongo.ChapterOnModerationPages{
		ChapterOnModerationID: newChapter.ID,
		CreatorID:             claims.ID,
		Pages:                 pages,
	}

	if _, err := h.ChaptersPages.InsertOne(c.Request.Context(), chapterPages); err != nil {
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
