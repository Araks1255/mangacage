package chapters

import (
	"context"
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
	"go.mongodb.org/mongo-driver/mongo"
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
		code, err := parseCreateChapterError(err)
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	err = insertChapterOnModerationPages(c.Request.Context(), h.ChaptersPages, requestBody.Pages, chapter.ID, claims.ID)
	if err != nil {
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

		case "volume":
			volume, err := strconv.ParseUint(string(data), 10, 64)
			if err != nil {
				return 400, err
			}
			body.Volume = uint(volume)

		case "titleId":
			id, err := strconv.ParseUint(string(data), 10, 64)
			if err != nil {
				return 400, err
			}
			idUint := uint(id)
			body.TitleID = &idUint

		case "titleOnModerationId":
			id, err := strconv.ParseUint(string(data), 10, 64)
			if err != nil {
				return 400, err
			}
			idUint := uint(id)
			body.TitleOnModerationID = &idUint

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

	if (parsedBody.TitleID != nil && parsedBody.TitleOnModerationID != nil) || (parsedBody.TitleID == nil && parsedBody.TitleOnModerationID == nil) {
		return 400, errors.New("ожидается один id тайтла")
	}

	if parsedBody.ID != nil {
		isOwner, err := helpers.CheckEntityOnModerationOwnership(db, "titles", *parsedBody.ID, userID)
		if err != nil {
			return 500, err
		}

		if !isOwner {
			return 403, errors.New("изменять заявку на модерацию может только её создатель")
		}
	}

	if parsedBody.TitleID != nil {
		var check struct {
			UserTeamID    *uint
			ChapterExists bool
		}

		err = db.Raw(
			`SELECT
				(
					SELECT tt.team_id FROM title_teams AS tt
					INNER JOIN users AS u ON u.team_id = tt.team_id
					WHERE tt.title_id = ? AND u.id = ?
				) AS user_team_id,
				EXISTS(
					SELECT 1 FROM chapters
					WHERE lower(name) = lower(?)
					AND title_id = ?
					AND volume = ?
					AND team_id = (SELECT team_id FROM users WHERE id = ?)
				) AS chapter_exists`,
			parsedBody.TitleID, userID, parsedBody.Name, parsedBody.TitleID, parsedBody.Volume, userID,
		).Scan(&check).Error

		if err != nil {
			return 500, err
		}

		if check.UserTeamID == nil {
			return 404, errors.New("тайтл не найден среди переводимых вашей командой")
		}
		if check.ChapterExists {
			return 409, errors.New("глава с таким названием уже выложена вашей командой в этом томе этого тайтла")
		}

		parsedBody.TeamID = *check.UserTeamID
	}

	if parsedBody.TitleOnModerationID != nil {
		var check struct {
			TitleOnModerationNew bool
			UserTeamID           *uint
		}

		err = db.Raw(
			`SELECT
				EXISTS(SELECT 1 FROM titles_on_moderation WHERE existing_id IS NULL AND id = ?) AS title_on_moderation_new,
				(SELECT team_id FROM users WHERE id = ?) AS user_team_id`,
			parsedBody.TitleOnModerationID, userID,
		).Scan(&check).Error

		if err != nil {
			return 500, err
		}

		if !check.TitleOnModerationNew {
			return 409, errors.New("не пытайтесь добавить том в изменения тайтла. просто добавьте их в уже существующий тайтл")
		}

		parsedBody.TeamID = *check.UserTeamID
	}

	return 0, nil
}

func parseCreateChapterError(err error) (code int, parsedError error) {
	if dbErrors.IsUniqueViolation(err, constraints.UniqChapterOnModerationVolumeTitleTeam) {
		return 409, errors.New("глава с таким названием и номером тома уже выложена вашей командой в этом тайтле")
	}

	if dbErrors.IsUniqueViolation(err, constraints.UniqChapterOnModerationVolumeTitleOnModeration) {
		return 409, errors.New("глава с таким названием и номером тома уже ожидает модерации в этом тайтле на модерации")
	}

	return 500, err
}

func insertChapterOnModerationPages(ctx context.Context, collection *mongo.Collection, pages [][]byte, chapterOnModerationID, userID uint) error {
	filter := bson.M{"chapter_on_moderation_id": chapterOnModerationID}
	update := bson.M{"$set": bson.M{"pages": pages, "creator_id": userID}}
	opts := options.Update().SetUpsert(true)

	if _, err := collection.UpdateOne(ctx, filter, update, opts); err != nil {
		return err
	}

	return nil
}
