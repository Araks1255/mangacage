package chapters

import (
	"errors"
	"log"
	"strconv"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/Araks1255/mangacage/pkg/common/utils"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	"github.com/Araks1255/mangacage/pkg/handlers/helpers"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h handler) EditChapter(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	chapter, err := mapRequestBodyToChapterOnModeration(c.Param, c.ShouldBindJSON, claims.ID)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	code, err := checkEditChapterConflicts(tx, *chapter)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	err = tx.Clauses(helpers.OnExistingIDConflictClause).Create(&chapter).Error

	if err != nil {
		code, err := parseChapterEditError(err)
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "изменения главы успешно отправлены на модерацию"})

	if _, err := h.NotificationsClient.NotifyAboutChapterOnModeration(c.Request.Context(), &pb.ChapterOnModeration{ID: uint64(*chapter.ExistingID), New: false}); err != nil {
		log.Println(err)
	}
}

func mapRequestBodyToChapterOnModeration(paramFn func(string) string, bindFn func(any) error, userID uint) (*models.ChapterOnModeration, error) {
	chapterID, err := strconv.ParseUint(paramFn("id"), 10, 64)
	if err != nil {
		return nil, err
	}

	var body dto.EditChapterDTO

	if err := bindFn(body); err != nil {
		return nil, err
	}

	ok, err := utils.HasAnyNonEmptyFields(body)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, errors.New("необходим как минимум 1 изменяемый параметр")
	}

	res := body.ToChapterOnModeration(userID, uint(chapterID))

	return &res, nil
}

func checkEditChapterConflicts(db *gorm.DB, chapter models.ChapterOnModeration) (code int, err error) {
	query := `SELECT
				EXISTS(
					SELECT 1 FROM chapters AS c
					INNER JOIN title_teams AS tt ON tt.title_id = c.title_id
					INNER JOIN users AS u ON tt.team_id = u.team_id
					WHERE c.id = ? AND u.id = ?
				) AS chapter_exists`

	if chapter.Name != nil {
		query += `,EXISTS(
			SELECT 1 FROM chapters
			WHERE lower(name) = lower(?)
			AND volume = (SELECT volume FROM chapters WHERE id = ?)
			AND team_id = (SELECT team_id FROM users WHERE id = ?)
		) AS chapter_with_the_same_name_exists`
	}

	if chapter.Volume != nil {
		query += `,EXISTS(
			SELECT 1 FROM chapters
			WHERE lower(name) = lower(?)
			AND volume = ?
			AND team_id = (SELECT team_id FROM users WHERE id = ?)
		) AS chapter_with_the_same_name_and_volume_exists`
	}

	if chapter.Name == nil && chapter.Volume == nil {
		var chapterExists bool
		if err := db.Raw(query, chapter.ExistingID, chapter.CreatorID).Scan(&chapterExists).Error; err != nil {
			return 500, err
		}
		if !chapterExists {
			return 404, errors.New("глава не найдена среди глав тайтлов переводимых вашей командой")
		}
	}

	if chapter.Name != nil && chapter.Volume == nil {
		var check struct {
			ChapterExists                bool
			ChapterWithTheSameNameExists bool
		}

		err := db.Raw(
			query,
			chapter.ExistingID, chapter.CreatorID,
			chapter.Name, chapter.ExistingID, chapter.CreatorID,
		).Scan(&check).Error

		if err != nil {
			return 500, err
		}

		if !check.ChapterExists {
			return 404, errors.New("глава не найдена среди глав тайтлов переводимых вашей командой")
		}
		if check.ChapterWithTheSameNameExists {
			return 409, errors.New("глава с таким названием и номером тома уже выложена вашей командой в этом тайтле")
		}
	}

	if chapter.Name != nil && chapter.Volume != nil {
		var check struct {
			ChapterExists                         bool
			ChapterWithTheSameNameExists          bool
			ChapterWithTheSameNameAndVolumeExists bool
		}

		err := db.Raw(
			query,
			chapter.ExistingID, chapter.CreatorID,
			chapter.Name, chapter.ExistingID, chapter.CreatorID,
			chapter.Name, chapter.Volume, chapter.CreatorID,
		).Scan(&check).Error

		if err != nil {
			return 500, err
		}

		if !check.ChapterExists {
			return 404, errors.New("глава не найдена среди глав тайтлов переводимых вашей командой")
		}
		if check.ChapterWithTheSameNameExists {
			return 409, errors.New("глава с таким названием и номером тома уже выложена вашей командой в этом тайтле")
		}
		if check.ChapterWithTheSameNameAndVolumeExists {
			return 409, errors.New("глава с таким названием и номером тома уже выложена вашей командой в этом тайтле")
		}
	}

	return 0, nil
}

func parseChapterEditError(err error) (code int, parsedError error) {
	if dbErrors.IsUniqueViolation(err, constraints.UniqChapterOnModerationVolumeTitleTeam) {
		return 409, errors.New("глава с таким названием и номером тома, созданная вашей командой, уже ожидает модерации в этом тайтле")
	}

	if dbErrors.IsUniqueViolation(err, constraints.UniqChapterOnModerationVolumeTitleOnModeration) {
		return 409, errors.New("глава с таким названием и номером тома уже ожидает модерации в этом тайтле на модерации")
	}

	return 500, err
}
