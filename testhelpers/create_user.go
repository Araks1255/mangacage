package testhelpers

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type CreateUserOptions struct {
	TeamID         uint
	TgUserID       int64
	Roles          []string
	Visible        bool
	ProfilePicture []byte
	PathToMediaDir string
}

func CreateUser(db *gorm.DB, opts ...CreateUserOptions) (uint, error) {
	if len(opts) > 1 {
		return 0, errors.New("Объект опций должен быть один")
	}

	tx := db.Begin()
	defer utils.RollbackOnPanic(tx)
	defer tx.Rollback()

	user := models.User{UserName: uuid.New().String()}

	if len(opts) != 0 {
		if opts[0].TeamID != 0 {
			user.TeamID = &opts[0].TeamID
		}
		if opts[0].TgUserID != 0 {
			user.TgUserID = &opts[0].TgUserID
		}
		user.Visible = opts[0].Visible

	}

	if result := tx.Create(&user); result.Error != nil {
		return 0, result.Error
	}

	if len(opts) == 0 {
		tx.Commit()
		return user.ID, nil
	}

	if len(opts[0].Roles) != 0 {
		if result := tx.Exec(
			`INSERT INTO user_roles (user_id, role_id)
			SELECT ?, roles.id FROM roles
			JOIN UNNEST(?::TEXT[]) AS role_name ON roles.name = role_name`,
			user.ID, pq.Array(opts[0].Roles),
		); result.Error != nil {
			return 0, result.Error
		}
	}

	if len(opts[0].ProfilePicture) != 0 {
		if opts[0].PathToMediaDir == "" {
			return 0, errors.New("не передана директория для сохранения медиафайлов")
		}

		path := fmt.Sprintf("%s/users/%d.jpg", opts[0].PathToMediaDir, user.ID)

		user.ProfilePicturePath = &path

		if err := os.MkdirAll(filepath.Dir(*user.ProfilePicturePath), 0644); err != nil {
			return 0, err
		}

		if err := os.WriteFile(*user.ProfilePicturePath, opts[0].ProfilePicture, 0644); err != nil {
			return 0, err
		}

		if err := tx.Updates(&user).Error; err != nil {
			return 0, err
		}
	}

	tx.Commit()

	return user.ID, nil
}
