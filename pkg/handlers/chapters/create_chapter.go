package chapters

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"strconv"
	"strings"
	"sync"

	"image"
	_ "image/jpeg"
	_ "image/png"

	"github.com/Araks1255/mangacage/pkg/auth"
	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	"github.com/Araks1255/mangacage/pkg/handlers/helpers"
	"github.com/Araks1255/mangacage_protos/gen/enums"
	pb "github.com/Araks1255/mangacage_protos/gen/site_notifications"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type counterReader struct {
	reader io.Reader
	count  int64
}

func (cr *counterReader) Read(p []byte) (n int, err error) {
	n, err = cr.reader.Read(p)
	cr.count += int64(n)
	return n, err
}

func (h handler) CreateChapter(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	reader, code, err := getRequestMultipartReader(c.Request.Header.Get, c.Request.MultipartReader)
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

	requestBody, tmpuuid, pages, code, err := parseCreateChapterBodyAndCreatePages(c.Request.Context(), reader, h.PathToMediaDir)

	if err != nil {
		if code == 500 {
			log.Println(err)
		}

		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})

		if err = removeTempPagesFolder(h.PathToMediaDir, tmpuuid); err != nil {
			log.Printf("не удалось удалить страницы главы на ошибке\nuuid главы: %s\nошибка: %s", tmpuuid, err.Error())
		}

		return
	}

	code, err = checkCreateChapterConflicts(tx, requestBody, claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	chapter := requestBody.ToChapterOnModeration(claims.ID)

	err = helpers.UpsertEntityOnModeration(tx, &chapter, chapter.ID)

	if err != nil {
		code, err := parseCreateChapterError(err)
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	errChan := make(chan error, 2)
	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		errChan <- createPagesMeta(tx, pages, h.PathToMediaDir, chapter.ID)
	}()
	go func() {
		defer wg.Done()
		errChan <- replaceChapterUUIDWithID(h.PathToMediaDir, tmpuuid, chapter.ID)
	}()

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})

			if err := removePagesFolder(h.PathToMediaDir, chapter.ID, tmpuuid); err != nil {
				log.Printf("не удалось удалить страницы главы на ошибке\nid главы: %d\nвременный uuid: %s\nошибка: %s", chapter.ID, tmpuuid, err.Error())
			}

			return
		}
	}

	if err := addPagesDirectoryPathToChapterOnModeration(tx, chapter.ID, h.PathToMediaDir); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "глава успешно отправлена на модерацию", "id": chapter.ID})

	if _, err := h.NotificationsClient.NotifyAboutNewModerationRequest(
		c.Request.Context(),
		&pb.ModerationRequest{
			EntityOnModeration: enums.EntityOnModeration_ENTITY_ON_MODERATION_CHAPTER,
			ID:                 uint64(chapter.ID),
		},
	); err != nil {
		log.Println(err)
	}

	if requestBody.NeedsCompression {
		h.ChaptersPagesCompressor.EnqueueChapterOnModerationID(chapter.ID)
	}
}

func getRequestMultipartReader(getHeaderFn func(string) string, getReaderFn func() (*multipart.Reader, error)) (reader *multipart.Reader, code int, err error) {
	contentType := getHeaderFn("Content-Type")

	if !strings.HasPrefix(contentType, "multipart/form-data") {
		return nil, 400, errors.New("тело запроса должно иметь тип multipart/form-data")
	}

	reader, err = getReaderFn()

	if err != nil {
		return nil, 500, err
	}

	return reader, 0, nil
}

func parseCreateChapterBodyAndCreatePages(
	ctx context.Context,
	r *multipart.Reader,
	pathToMediaDir string,
) (
	requestBody *dto.CreateChapterDTO,
	tmpuuid string,
	pages []models.Page,
	code int,
	err error,
) {
	var (
		body                                    dto.CreateChapterDTO
		usedNumbers                             = make(map[uint]struct{}, 40)
		maxPageQuality                          = 2.5 // Максимальное качество одной страницы. Измеряется в байтах на пиксель. Эквивалентно странице 1080:1920 весом 5мб
		maxAveragePageQuality                   = 1.7 // Максимальное среднее качество страниц. Эквивалентно странице 1080:1920 весом 3.5мб
		maxAveragePageQualityWithoutCompression = 0.5 // Максимальное среднее качество страниц без сжатия. Эквивалентно странице 1080:1920 весом 1мб
	)
	pages = make([]models.Page, 0, 40)
	tmpuuid = uuid.NewString()
Loop:
	for {
		select {
		case <-ctx.Done():
			return
		default:
			part, err := r.NextPart()

			if err == io.EOF {
				break Loop
			}
			if err != nil {
				return nil, "", nil, 500, err
			}

			limitReader := io.LimitReader(part, 135<<20+1) // Больше 100 все равно не пропустит реверс прокси, но это как самая крайняя мера предосторожности. А ограничения нет, ведь в теории можно и всю главу в одно изображение склеить

			if strings.HasPrefix(part.FormName(), "page") {
				fieldName := []rune(part.FormName())

				body.NumberOfPages++

				if body.NumberOfPages == 1500 {
					return nil, "", nil, 400, errors.New("превышен лимит страниц (1500)")
				}

				number, err := strconv.ParseUint(string(fieldName[4:]), 10, 32)
				if err != nil {
					return nil, "", nil, 400, errors.New("неверный формат имени поля страницы. именование должно быть вида page1, page2...")
				}

				if _, ok := usedNumbers[uint(number)]; ok {
					return nil, "", nil, 400, errors.New("номера страниц не должны повторяться")
				}

				if part.FileName() == "" {
					return nil, "", nil, 400, errors.New("в поле страницы отправлен не файл")
				}

				if err := os.MkdirAll(fmt.Sprintf("%s/chapters_on_moderation/%s", pathToMediaDir, tmpuuid), 0755); err != nil {
					return nil, "", nil, 500, err
				}

				tempFile, err := os.CreateTemp(fmt.Sprintf("%s/chapters_on_moderation/%s/", pathToMediaDir, tmpuuid), "*.tmp")
				if err != nil {
					return nil, "", nil, 500, err
				}

				counterLimitReader := counterReader{reader: limitReader}
				counterLimitTeeReader := io.TeeReader(&counterLimitReader, tempFile)

				config, format, err := image.DecodeConfig(counterLimitTeeReader)
				if err != nil {
					tempFile.Close()
					return nil, "", nil, 400, errors.New("в поле страницы отправлено не фото (разрешенные форматы: jpg, png)")
				}

				pageResolution := float64(config.Width * config.Height)

				aspectRatio := float64(config.Height) / float64(config.Width)
				isWebtoonPage := aspectRatio > 2.4

				if isWebtoonPage {
					body.WebtoonMode = true // Для больших изображений, содержащих несколько страниц (как в манхве) включается webtoon mode
				}

				if _, err := io.Copy(io.Discard, counterLimitTeeReader); err != nil {
					tempFile.Close()
					return nil, "", nil, 500, err
				}

				body.PagesSize += counterLimitReader.count
				body.PagesResolution += int64(pageResolution)

				if body.PagesSize > 135<<20 {
					tempFile.Close()
					return nil, "", nil, 400, errors.New("превышен максимальный размер главы (135мб)")
				}

				if float64(counterLimitReader.count)/pageResolution > maxPageQuality {
					tempFile.Close()
					return nil, "", nil, 400, fmt.Errorf("превышено максимальное качество страницы (%fбайт/пиксель)", maxPageQuality)
				}

				tempFile.Close()

				if err := os.Rename(tempFile.Name(), fmt.Sprintf("%s/chapters_on_moderation/%s/%d.%s", pathToMediaDir, tmpuuid, number, format)); err != nil {
					return nil, "", nil, 500, err
				}

				pages = append(pages, models.Page{Number: uint(number), Format: format})

				usedNumbers[uint(number)] = struct{}{}

				continue
			}

			data, err := io.ReadAll(limitReader)
			if err != nil {
				return nil, "", nil, 500, err
			}

			switch part.FormName() {
			case "name":
				body.Name = string(data)
				if len([]rune(body.Name)) > 80 {
					return nil, "", nil, 400, errors.New("превышена максимальная длина названия (80 символов)")
				}

			case "id":
				id, err := strconv.ParseUint(string(data), 10, 64)
				if err != nil {
					return nil, "", nil, 400, err
				}
				idUint := uint(id)
				body.ID = &idUint

			case "volume":
				volume, err := strconv.ParseUint(string(data), 10, 64)
				if err != nil {
					return nil, "", nil, 400, err
				}
				body.Volume = uint(volume)

			case "titleId":
				id, err := strconv.ParseUint(string(data), 10, 64)
				if err != nil {
					return nil, "", nil, 400, err
				}
				idUint := uint(id)
				body.TitleID = &idUint

			case "titleOnModerationId":
				id, err := strconv.ParseUint(string(data), 10, 64)
				if err != nil {
					return nil, "", nil, 400, err
				}
				idUint := uint(id)
				body.TitleOnModerationID = &idUint

			case "description":
				description := string(data)
				if len([]rune(description)) > 100 {
					return nil, "", nil, 400, errors.New("превышена максимальная длина описания (100 символов)")
				}
				body.Description = &description

			case "disableCompression":
				body.DisableCompression, err = strconv.ParseBool(string(data))
				if err != nil {
					return nil, "", nil, 400, err
				}

			case "webtoonMode": // Явное включение webtoon режима
				body.WebtoonMode, err = strconv.ParseBool(string(data))
				if err != nil {
					return nil, "", nil, 400, err
				}

			default:
				return nil, "", nil, 400, errors.New("отправлено неизвестное поле. поля страниц именуются в формате page1, page2...")
			}
		}
	}

	averagePageQuality := float64(body.PagesSize) / float64(body.PagesResolution)

	if body.DisableCompression && averagePageQuality > maxAveragePageQualityWithoutCompression {
		return nil, "", nil, 400, fmt.Errorf("превышено максимальное среднее качество страницы без сжатия (%fбайт/пиксель)", maxAveragePageQualityWithoutCompression)
	}

	if averagePageQuality > maxAveragePageQuality {
		return nil, "", nil, 400, fmt.Errorf("превышено максимальное среднее качество страницы (%fбайт/пиксель)", maxAveragePageQuality)
	}

	body.NeedsCompression = !body.DisableCompression && averagePageQuality > maxAveragePageQualityWithoutCompression

	return &body, tmpuuid, pages, 0, nil
}

func removeTempPagesFolder(pathToMediaDir, tmpuuid string) error {
	return os.RemoveAll(fmt.Sprintf("%s/chapters_on_moderation/%s", pathToMediaDir, tmpuuid))
}

func removePagesFolder(pathToMediaDir string, chapterOnModerationID uint, tmpuuid string) error {
	err := os.RemoveAll(fmt.Sprintf("%s/chapters_on_moderation/%d", pathToMediaDir, chapterOnModerationID))

	if err != nil {
		if os.IsNotExist(err) {
			return os.RemoveAll(fmt.Sprintf("%s/chapters_on_moderation/%s", pathToMediaDir, tmpuuid))
		}
		return err
	}

	return nil
}

func checkCreateChapterConflicts(db *gorm.DB, parsedBody *dto.CreateChapterDTO, userID uint) (code int, err error) {
	if parsedBody.Name == "" || parsedBody.NumberOfPages == 0 || parsedBody.Volume == 0 || (parsedBody.TitleID == nil && parsedBody.TitleOnModerationID == nil) {
		return 400, errors.New("в запросе недостаточно данных")
	}

	if (parsedBody.TitleID != nil && parsedBody.TitleOnModerationID != nil) || (parsedBody.TitleID == nil && parsedBody.TitleOnModerationID == nil) {
		return 400, errors.New("ожидается один id тайтла")
	}

	if parsedBody.ID != nil {
		isOwner, err := helpers.CheckEntityOnModerationOwnership(db, "chapters", *parsedBody.ID, userID)
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
			return 409, errors.New("не пытайтесь добавить главу в изменения тайтла. просто добавьте ее в уже существующий тайтл")
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

func addPagesDirectoryPathToChapterOnModeration(db *gorm.DB, id uint, pathToMediaDir string) error {
	return db.Exec(
		"UPDATE chapters_on_moderation SET pages_dir_path = ? WHERE id = ?",
		fmt.Sprintf("%s/chapters_on_moderation/%d", pathToMediaDir, id), id,
	).Error
}

func replaceChapterUUIDWithID(pathToMediaDir, tmpuuid string, id uint) error {
	oldPath := fmt.Sprintf("%s/chapters_on_moderation/%s", pathToMediaDir, tmpuuid)
	newPath := fmt.Sprintf("%s/chapters_on_moderation/%d", pathToMediaDir, id)
	return os.Rename(oldPath, newPath)
}

func createPagesMeta(db *gorm.DB, pages []models.Page, pathToMediaDir string, chapterOnModerationID uint) error {
	for i := 0; i < len(pages); i++ {
		path := fmt.Sprintf("%s/chapters_on_moderation/%d/%d.%s", pathToMediaDir, chapterOnModerationID, pages[i].Number, pages[i].Format)
		pages[i].ChapterOnModerationID = &chapterOnModerationID
		pages[i].Path = path
	}
	return db.Create(&pages).Error
}
