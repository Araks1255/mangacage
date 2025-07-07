package chapters

import (
	"errors"
	"io"
	"log"
	"mime/multipart"
	"strconv"
	"strings"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	"github.com/Araks1255/mangacage/pkg/handlers/helpers"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gorm.io/gorm"
)

func (h handler) CreateChapter(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

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

	var requestBody dto.CreateChapterDTO

	code, err := parseCreateChapterBody(reader, &requestBody)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	code, err = checkCreateChapterConflicts(h.DB, &requestBody, claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	chapter := requestBody.ToChapterOnModeration(claims.ID)

	err = tx.Clauses(helpers.OnIDConflictClause).Create(&chapter).Error

	if err != nil {
		if dbErrors.IsUniqueViolation(err, constraints.UniqChapterOnModerationVolume) {
			c.AbortWithStatusJSON(409, gin.H{"error": "глава с таким названием уже ожидает модерации в этом томе"})
			return
		}
		if dbErrors.IsUniqueViolation(err, constraints.UniqChapterOnModerationVolumeOnModeration) {
			c.AbortWithStatusJSON(409, gin.H{"error": "глава с таким названием уже ожидает модерации в этом томе на модерации"})
			return
		}

		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	filter := bson.M{"chapter_on_moderation_id": chapter.ID}
	update := bson.M{"$set": bson.M{"pages": requestBody.Pages, "creator_id": claims.ID}}
	opts := options.Update().SetUpsert(true)

	if _, err := h.ChaptersPages.UpdateOne(c.Request.Context(), filter, update, opts); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "глава успешно отправлена на модерацию"})

	if _, err := h.NotificationsClient.NotifyAboutChapterOnModeration(c.Request.Context(), &pb.ChapterOnModeration{ID: uint64(chapter.ID), New: true}); err != nil {
		log.Println(err)
	}
}

func parseCreateChapterBody(r *multipart.Reader, body *dto.CreateChapterDTO) (code int, err error) {
	body.Pages = make([][]byte, 0, 45)
	for {
		part, err := r.NextPart()

		if err == io.EOF {
			break
		}
		if err != nil {
			return 500, err
		}

		data, err := io.ReadAll(part)
		if err != nil {
			return 500, err
		}

		switch part.FormName() {
		case "name":
			body.Name = string(data)

		case "id":
			id, err := strconv.ParseUint(string(data), 10, 64)
			if err != nil {
				return 400, err
			}
			idUint := uint(id)
			body.ID = &idUint

		case "volumeId":
			id, err := strconv.ParseUint(string(data), 10, 64)
			if err != nil {
				return 400, err
			}
			idUint := uint(id)
			body.VolumeID = &idUint

		case "volumeOnModerationId":
			id, err := strconv.ParseUint(string(data), 10, 64)
			if err != nil {
				return 400, err
			}
			idUint := uint(id)
			body.VolumeOnModerationID = &idUint

		case "description":
			description := string(data)
			body.Description = &description

		case "pages":
			if len(data) > 1<<20 {
				return 400, errors.New("превышен максимальный размер страницы (1мб)")
			}

			body.Pages = append(body.Pages, data)

			if len(body.Pages) > 230 {
				return 400, errors.New("превышено максимальное количество страниц (230)")
			}
		}
	}

	return 0, nil
}

func checkCreateChapterConflicts(db *gorm.DB, parsedBody *dto.CreateChapterDTO, userID uint) (code int, err error) {
	if parsedBody.Name == "" || len(parsedBody.Pages) == 0 {
		return 400, errors.New("в запросе недостаточно данных")
	}

	if (parsedBody.VolumeID != nil && parsedBody.VolumeOnModerationID != nil) || (parsedBody.VolumeID == nil && parsedBody.VolumeOnModerationID == nil) {
		return 400, errors.New("ожидается один id тома")
	}

	if parsedBody.VolumeID != nil {
		var check struct {
			ChapterExists bool
			UserTeamID    *uint
		}

		err := db.Raw(
			`SELECT
				EXISTS(SELECT 1 FROM chapters WHERE lower(name) = lower(?) AND volume_id = ?) AS chapter_exists,
				(
					SELECT
						tt.team_id
					FROM
						title_teams AS tt
						INNER JOIN volumes AS v ON v.title_id = tt.title_id
						INNER JOIN users AS u ON u.team_id = tt.team_id
					WHERE
						v.id = ? AND u.id = ?
				) AS user_team_id`,
			parsedBody.Name, *parsedBody.VolumeID, *parsedBody.VolumeID, userID,
		).Scan(&check).Error

		if err != nil {
			return 500, err
		}

		if check.UserTeamID == nil {
			return 404, errors.New("том не найден среди переводимых вашей командрй")
		}
		if check.ChapterExists {
			return 409, errors.New("глава с таким названием уже существует в этом томе")
		}

		parsedBody.TeamID = *check.UserTeamID
	}

	if parsedBody.VolumeOnModerationID != nil {
		var userTeamID *uint

		err := db.Raw(
			`SELECT u.team_id FROM users AS u
			INNER JOIN volumes_on_moderation AS vom ON vom.creator_id = u.id
			WHERE vom.id = ? AND u.id = ?`,
			*parsedBody.VolumeOnModerationID, userID,
		).Scan(&userTeamID).Error

		if err != nil {
			return 500, err
		}

		if userTeamID == nil {
			return 404, errors.New("том на модерации не найден среди ваших заявок")
		}

		parsedBody.TeamID = *userTeamID
	}

	return 0, nil
}
