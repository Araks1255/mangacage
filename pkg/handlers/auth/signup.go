package auth

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/auth/utils"
	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	"github.com/Araks1255/mangacage/pkg/handlers/helpers"
	pb "github.com/Araks1255/mangacage_protos"

	"github.com/gin-gonic/gin"
)

func (h handler) Signup(c *gin.Context) {
	if _, err := c.Cookie("mangacage_token"); err == nil {
		c.AbortWithStatusJSON(409, gin.H{"error": "вы уже вошли в аккаунт"})
		return
	}

	var requestBody models.UserOnModerationDTO

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	exists, err := helpers.CheckEntityWithTheSameNameExistence(h.DB, "users", requestBody.UserName, nil, nil)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if exists {
		c.AbortWithStatusJSON(409, gin.H{"error": "пользователь с таким именем уже существует"})
		return
	}

	hash, err := utils.GenerateHashPassword(*requestBody.Password)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}
	requestBody.Password = &hash

	user := requestBody.ToUserOnModeration(nil)

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	err = h.DB.Create(&user).Error

	if err != nil {
		if dbErrors.IsUniqueViolation(err, constraints.UniqUserOnModerationUserName) { // Такие проверки всё же гораздо быстрее чем ручной SELECT перед вставкой. Ну, расчёт на то, что на другую бд мигрировать проект не будет и автосгенерированные названия ограничений не поменяются
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
