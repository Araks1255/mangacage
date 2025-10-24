package chapters_pages_compressor

import (
	"bytes"
	"image"
	"image/jpeg"
	_ "image/png"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Araks1255/mangacage/pkg/common/models"

	"github.com/nfnt/resize"
)

func (c *ChaptersPagesCompressor) Start() {
MainLoop:
	for chapteronModerationID := range c.chaptersOnModerationIDs {
		start := time.Now()

		var pagesMeta []models.Page

		if err := c.db.Raw("SELECT * FROM pages WHERE chapter_on_moderation_id = ?", chapteronModerationID).Scan(&pagesMeta).Error; err != nil {
			log.Printf(
				"ошибка при получении метаданных страниц главы на модерации\nid главы на модерации: %d\nошибка: %s",
				chapteronModerationID, err.Error(),
			)
		}

		if len(pagesMeta) == 0 {
			log.Printf("не найдено метаданных страниц главы на модерации\nid главы на модерации: %d", chapteronModerationID)
		}

		for i := 0; i < len(pagesMeta); i++ {
			info, err := os.Stat(pagesMeta[i].Path)

			if err != nil {
				if !os.IsNotExist(err) {
					log.Printf("не удалось открыть файл по ошибке: %s", err.Error())
					continue
				}

				oldPath := pagesMeta[i].Path

				err = c.db.Raw(
					"SELECT * FROM pages WHERE chapter_on_moderation_id = ? AND number = ?",
					chapteronModerationID, pagesMeta[i].Number,
				).Scan(&pagesMeta[i]).Error

				if err != nil {
					log.Printf(
						"ошибка при попытке получения новых метаданных ненайденного файла\nid главы на модерации: %d\nномер страницы: %d\nошибка: %s",
						chapteronModerationID, pagesMeta[i].Number, err.Error(),
					)
				}

				if pagesMeta[i].Path == "" {
					log.Printf(
						"новый путь к файлу не найден. возможно, глава была отклонена в процессе сжатия\nid главы на модерации: %d",
						chapteronModerationID,
					)
					continue MainLoop
				}

				info, err = os.Stat(pagesMeta[i].Path)
				if err != nil {
					log.Printf(
						"ошибка при повторном получении файла с новым путем\nid главы на модерации: %d\nid главы: %d\nстарый путь: %s\nновый путь: %s\nошибка: %s",
						chapteronModerationID, pagesMeta[i].ChapterID, oldPath, pagesMeta[i].Path, err.Error(),
					)
					continue MainLoop
				}
			}

			if info.Size() < 600*1000 {
				continue
			}

			data, err := os.ReadFile(pagesMeta[i].Path)

			if err != nil {
				log.Printf(
					"не удалось прочитать файл по ошибке\nid главы на модерации: %d\nпуть к файлу: %s\nошибка: %s",
					chapteronModerationID, pagesMeta[i].Path, err.Error(),
				)
				continue
			}

			img, _, err := image.Decode(bytes.NewReader(data))

			if err != nil {
				log.Printf(
					"ошибка при декодировании файла в фото\nid главы на модерации: %d\nпуть к файлу: %s\nошибка: %s",
					chapteronModerationID, pagesMeta[i].Path, err.Error(),
				)
			}

			if img.Bounds().Dy() > 1080 {
				img = resize.Resize(0, 1080, img, resize.Bilinear)
			}

			var buf bytes.Buffer

			if err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 85}); err != nil {
				log.Printf(
					"ошибка при сжатии файла\nid главы на модерации: %d\nпуть к файлу: %s\n ошибка: %s",
					chapteronModerationID, pagesMeta[i].Path, err.Error(),
				)
				continue
			}

			if err := os.WriteFile(pagesMeta[i].Path, buf.Bytes(), 0755); err != nil {
				log.Printf(
					"не удалось записать в файл сжатое содержимое\nid главы на модерации: %d\nпуть к файлу: %s\nошибка: %s",
					chapteronModerationID, pagesMeta[i].Path, err.Error(),
				)
				continue
			}

			if pagesMeta[i].Format != "jpg" && pagesMeta[i].Format != "jpeg" {

				pagesMeta[i].Format = "jpeg"

				if pagesMeta[i].Path, err = changeFileExtension(pagesMeta[i].Path, ".jpeg"); err != nil {
					log.Printf("ошибка при изменении расширения файла после сжатия\nпуть: %s\nошибка: %s", pagesMeta[i].Path, err.Error())
					continue
				}

				if err := c.db.Updates(&pagesMeta[i]).Error; err != nil {
					log.Printf(
						"не удалось обновить метаинформацию о странице\nid страницы: %d\nid главы на модерации: %d\nномер страницы: %d",
						pagesMeta[i].ID, chapteronModerationID, pagesMeta[i].Number,
					)
					continue
				}
			}
		}

		log.Printf("выполнено за %s", time.Since(start))
	}
}

func changeFileExtension(oldPath, newExtension string) (newPath string, err error) {
	ext := filepath.Ext(oldPath)
	newPath = strings.TrimSuffix(oldPath, ext) + newExtension
	return newPath, os.Rename(oldPath, newPath)
}
