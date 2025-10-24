package moderation

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

func (h handler) CancelAppealForTitleModeration(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	titleOnModerationID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": "указан невалидный id тайтла на модерации"})
		return
	}

	coverPath, chaptersPagesDirectoriesPathes, code, err := deleteTitleOnModeration(h.DB, uint(titleOnModerationID), claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": "заявка на модерацию тайтла успешно отменена"})

	if coverPath != nil || len(chaptersPagesDirectoriesPathes) != 0 {
		errChan := make(chan error, 2)
		var wg sync.WaitGroup

		if coverPath != nil {
			wg.Add(1)
			go func() {
				defer wg.Done()
				errChan <- deleteTitleOnModerationCover(*coverPath)
			}()
		}
		if len(chaptersPagesDirectoriesPathes) != 0 {
			wg.Add(1)
			go func() {
				defer wg.Done()
				errChan <- deleteChaptersOnModerationPages(chaptersPagesDirectoriesPathes)
			}()
		}

		wg.Wait()
		close(errChan)

		for err := range errChan {
			if err != nil {
				log.Println(err)
			}
		}
	}
}

func deleteTitleOnModeration(db *gorm.DB, id, creatorID uint) (coverPath *string, chaptersPagesDirectoriesPathes []string, code int, err error) {
	var res struct {
		ID                             *uint
		CoverPath                      *string
		ChaptersPagesDirectoriesPathes pq.StringArray `gorm:"type:TEXT[]"`
	}

	query :=
		`WITH deleting_chapters AS (
			SELECT ARRAY(
				SELECT
					pages_dir_path
				FROM
					chapters_on_moderation
				WHERE
					title_on_moderation_id = ?
			) AS chapters_pages_directories_pathes
		)
		DELETE FROM
			titles_on_moderation
		WHERE
			id = ? AND creator_id = ?
		RETURNING
			id, cover_path, (SELECT chapters_pages_directories_pathes FROM deleting_chapters)::TEXT[]`

	if err := db.Raw(query, id, id, creatorID).Scan(&res).Error; err != nil {
		return nil, nil, 500, err
	}

	if res.ID == nil {
		return nil, nil, 404, errors.New("заявка на модерацию тайтла не найдена среди отправленных вами")
	}

	return res.CoverPath, res.ChaptersPagesDirectoriesPathes, 0, nil
}

func deleteTitleOnModerationCover(path string) error {
	if err := os.Remove(path); err != nil {
		return fmt.Errorf("ошибка при удалении обложки тайтла на модерации\nпуть: %s\nошибка: %s", path, err.Error())
	}
	return nil
}

func deleteChaptersOnModerationPages(dirs []string) error {
	if len(dirs) == 0 {
		return nil
	}

	errChan := make(chan error, len(dirs))
	var wg sync.WaitGroup

	wg.Add(len(dirs))

	for i := 0; i < len(dirs); i++ {
		go func(path string) {
			defer wg.Done()
			errChan <- os.RemoveAll(path)
		}(dirs[i])
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return fmt.Errorf(
				"ошибка при удалении страниц глав на модерации удаляемого тайтла на модераци\nдиректории страниц: %v\nошибка: %s",
				dirs, err.Error(),
			)
		}
	}

	return nil
}
