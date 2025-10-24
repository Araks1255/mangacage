package teams

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"log"
	"mime/multipart"
	"os"

	"gorm.io/gorm"
)

func CreateTeamOnModerationCover(db *gorm.DB, pathToMediaDir string, id uint, cover *multipart.FileHeader) (code int, err error) {
	if cover == nil {
		return 0, err
	}

	file, err := cover.Open()
	if err != nil {
		return 500, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return 500, err
	}

	_, format, err := image.DecodeConfig(bytes.NewReader(data))
	if err != nil {
		return 400, errors.New("ошибка при декодировании файла. скорее всего, было отправлено не фото")
	}

	path := fmt.Sprintf("%s/teams_on_moderation/%d.%s", pathToMediaDir, id, format)

	var oldPath *string

	if err := db.Raw("SELECT cover_path FROM teams_on_moderation WHERE id = ?", id).Scan(&oldPath).Error; err != nil {
		log.Printf("ошибка при получении старого пути к обложе команды на модерации\nid команды на модерации: %d\nошибка: %s", id, err.Error())
		return 500, err // Я бы не хотел возвращать здесь ошибку, так как запрос не такой уж и важный, но postgres всё равно блокирует транзакцию при ошибке
	}

	result := db.Exec("UPDATE teams_on_moderation SET cover_path = ? WHERE id = ?", path, id)

	if result.Error != nil {
		return 500, result.Error
	}

	if result.RowsAffected == 0 {
		return 500, fmt.Errorf("не удалось добавить путь к обложке команды на модерации\nid команды на модерации: %d", id)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return 500, err
	}

	if oldPath != nil && *oldPath != path {
		if err := os.Remove(*oldPath); err != nil {
			log.Printf(
				"не удалось удалить старый файл с обложкой команды на модерации\nid команды на модерации: %d\nпуть: %s\nошибка: %s",
				id, *oldPath, err.Error(),
			)
		}
	}

	return 0, nil
}
