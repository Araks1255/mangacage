package testhelpers

import (
	"errors"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

type CreateUserOptions struct {
	TeamID uint
	Roles  []string
}

func CreateUser(db *gorm.DB, opts ...CreateUserOptions) (uint, error) {
	if len(opts) > 1 {
		return 0, errors.New("Объект опций должен быть один")
	}

	tx := db.Begin()
	defer func() { // Вынесу потом эту функцию куда-нибудь
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()
	defer tx.Rollback()

	user := models.User{
		UserName: uuid.New().String(),
	}

	if result := tx.Create(&user); result.Error != nil {
		return 0, result.Error
	}

	if len(opts) == 0 {
		tx.Commit()
		return user.ID, nil
	}

	if opts[0].TeamID != 0 {
		if result := tx.Exec("UPDATE users SET team_id = ? WHERE id = ?", opts[0].TeamID, user.ID); result.Error != nil {
			return 0, result.Error
		}
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

	tx.Commit()
	return user.ID, nil
}
