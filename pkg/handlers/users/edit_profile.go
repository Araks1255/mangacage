package users

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

	"github.com/Araks1255/mangacage/pkg/auth"
	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/Araks1255/mangacage/pkg/common/utils"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	"github.com/Araks1255/mangacage/pkg/handlers/helpers"
	"github.com/Araks1255/mangacage_protos/gen/enums"
	pb "github.com/Araks1255/mangacage_protos/gen/site_notifications"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"gorm.io/gorm"
)

func (h handler) EditProfile(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	var requestBody dto.EditUserDTO

	if err := c.ShouldBindWith(&requestBody, binding.FormMultipart); err != nil {
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	code, err := checkEditProfileConflicts(h.DB, requestBody)
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

	editedProfile := requestBody.ToUserOnModeration(claims.ID)

	err = helpers.UpsertEntityChanges(tx, editedProfile, *editedProfile.ExistingID)

	if err != nil {
		if dbErrors.IsUniqueViolation(err, constraints.UniqUserOnModerationUserName) {
			c.AbortWithStatusJSON(409, gin.H{"error": "пользователь с таким именем уже ожидает модерации"})
		} else {
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}
		return
	}

	if requestBody.ProfilePicture != nil {
		if code, err := createUserProfilePicture(tx, h.PathToMediaDir, editedProfile.ID, requestBody.ProfilePicture); err != nil {
			if code == 500 {
				log.Println(err)
			}
			c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
			return
		}
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "изменения профиля успешно отправлены на модерацию"})

	if _, err := h.NotificationsClient.NotifyAboutNewModerationRequest(
		c.Request.Context(),
		&pb.ModerationRequest{
			EntityOnModeration: enums.EntityOnModeration_ENTITY_ON_MODERATION_PROFILE_CHANGES,
			ID:                 uint64(editedProfile.ID),
		},
	); err != nil {
		log.Println(err)
	}
}

func checkEditProfileConflicts(db *gorm.DB, requestBody dto.EditUserDTO) (code int, err error) {
	ok, err := utils.HasAnyNonEmptyFields(&requestBody)
	if err != nil {
		return 500, err
	}

	if !ok {
		return 400, errors.New("необходим как минимум 1 изменяемый параметр")
	}

	if requestBody.UserName != nil {
		exists, err := helpers.CheckEntityWithTheSameNameExistence(db, "users", requestBody.UserName, nil, nil)
		if err != nil {
			return 500, err
		}

		if exists {
			return 409, errors.New("пользователь с таким именем уже существует")
		}
	}

	return 0, nil
}

func createUserProfilePicture(db *gorm.DB, pathToMediaDir string, id uint, profilePicture *multipart.FileHeader) (code int, err error) {
	if profilePicture == nil {
		return 0, nil
	}

	file, err := profilePicture.Open()
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

	path := fmt.Sprintf("%s/users_on_moderation/%d.%s", pathToMediaDir, id, format)

	var oldPath *string

	if err := db.Raw("SELECT profile_picture_path FROM users_on_moderation WHERE id = ?", id).Scan(&oldPath).Error; err != nil {
		log.Printf("ошибка при получении старого пути к аватарке изменений профиля\nid изменений профиля: %d\nошибка: %s", id, err.Error())
		return 500, err // Я бы не хотел возвращать здесь ошибку, так как запрос не такой уж и важный, но postgres всё равно блокирует транзакцию при ошибке
	}

	result := db.Exec("UPDATE users_on_moderation SET profile_picture_path = ? WHERE id = ?", path, id)

	if result.Error != nil {
		return 500, result.Error
	}

	if result.RowsAffected == 0 {
		return 500, errors.New("не удалось добавить путь к аватарке изменений профиля")
	}

	if err := os.WriteFile(path, data, 0755); err != nil {
		return 500, err
	}

	if oldPath != nil && *oldPath != path {
		if err := os.Remove(*oldPath); err != nil {
			log.Printf(
				"не удалось удалить старый файл с аватаркой изменений профиля\nid изменений профиля: %d\nпуть: %s\nошибка: %s",
				id, *oldPath, err.Error(),
			)
		}
	}

	return 0, nil
}
