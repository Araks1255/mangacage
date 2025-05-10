package auth

import (
	"database/sql"
	"log"

	"github.com/Araks1255/mangacage/pkg/auth/utils"
	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	pb "github.com/Araks1255/mangacage_protos"

	"github.com/gin-gonic/gin"
)

func (h handler) Signup(c *gin.Context) {
	if _, err := c.Cookie("mangacage_token"); err == nil {
		c.AbortWithStatusJSON(409, gin.H{"error": "вы уже вошли в аккаунт"})
		return
	}

	var requestBody struct {
		UserName      string `json:"userName" binding:"required"`
		Password      string `json:"password" binding:"required,min=8"`
		AboutYourself string `json:"aboutYourself"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	user := models.UserOnModeration{
		UserName:      sql.NullString{String: requestBody.UserName, Valid: true},
		AboutYourself: requestBody.AboutYourself,
	}

	var errHash error
	user.Password, errHash = utils.GenerateHashPassword(requestBody.Password)
	if errHash != nil {
		log.Println(errHash)
		c.AbortWithStatusJSON(500, gin.H{"error": errHash.Error()})
		return
	}

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var doesUserExist bool

	if err := tx.Raw("SELECT EXISTS(SELECT 1 FROM users WHERE lower(user_name) = lower(?))", requestBody.UserName).Error; err != nil { // Тут по-прежнему ручной SELECT, потому-что через индексы ограничение уникальности по двум таблицам не сделаешь (напрямую как минимум)
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if doesUserExist {
		c.AbortWithStatusJSON(409, gin.H{"error": "пользователь с таким именем уже существует"})
		return
	}

	err := h.DB.Create(&user).Error

	if err != nil {
		if dbErrors.IsUniqueViolation(err, constraints.UniUsersOnModerationUsername) { // Такие проверки всё же гораздо быстрее чем ручной SELECT перед вставкой. Ну, расчёт на то, что на другую бд мигрировать проект не будет и автосгенерированные названия ограничений не поменяются
			c.AbortWithStatusJSON(409, gin.H{"error": "пользователь с таким именем уже ожидает модерации"})
		} else {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "Ваш аккаунт успешно создан и ожидает верификации"})

	if _, err := h.NotificationsClient.NotifyAboutUserOnModeration(c.Request.Context(), &pb.User{ID: uint64(user.ID), New: true}); err != nil {
		log.Println(err)
	}
}
