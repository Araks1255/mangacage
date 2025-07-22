package users

import (
	"errors"
	"log"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/Araks1255/mangacage/pkg/common/utils"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func (h handler) EditProfileOnVerification(c *gin.Context) {
	claims := c.MustGet("claims").(*auth.Claims)

	user, code, err := mapEditProfileOnVerificationRequestodyToUser(c.ShouldBindJSON, claims.ID)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	code, err = updateProfileOnVerification(h.DB, *user)
	if err != nil {
		if code == 500 {
			log.Println(err)
		}
		c.AbortWithStatusJSON(code, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"success": "изменения профиля успешно применены"})
}

func mapEditProfileOnVerificationRequestodyToUser(bindFn func(any) error, userID uint) (user *models.User, code int, err error) {
	var requestBody dto.EditUserOnVerificationDTO

	if err = bindFn(&requestBody); err != nil {
		return nil, 400, err
	}

	ok, err := utils.HasAnyNonEmptyFields(&requestBody)
	if err != nil {
		return nil, 500, err
	}

	if !ok {
		return nil, 400, errors.New("необходим как минимум 1 изменяемый параметр")
	}

	res := requestBody.ToUser(userID)

	return &res, 0, nil
}

func updateProfileOnVerification(db *gorm.DB, user models.User) (code int, err error) {
	result := db.Table("users").Updates(&user).Where("!verificated")

	if result.Error != nil {
		return 500, result.Error
	}

	if result.RowsAffected == 0 {
		return 409, errors.New("ваш аккаунт уже прошел верификацию")
	}

	return 0, nil
}
