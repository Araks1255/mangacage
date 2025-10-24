package auth

import (
	"log"

	"github.com/Araks1255/mangacage/pkg/auth/utils"
	dbErrors "github.com/Araks1255/mangacage/pkg/common/db/errors"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"
	"github.com/Araks1255/mangacage/pkg/constants/postgres/constraints"
	pb "github.com/Araks1255/mangacage_protos/gen/site_notifications"

	"github.com/gin-gonic/gin"
)

func (h handler) Signup(c *gin.Context) {
	_, err := c.Cookie("mangacage_token")
	if err == nil {
		c.AbortWithStatusJSON(409, gin.H{"error": "вы уже вошли в аккаунт"})
		return
	}

	var requestBody dto.CreateUserDTO

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	requestBody.Password, err = utils.GenerateHashPassword(requestBody.Password)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	user := requestBody.ToUser()

	err = h.DB.Create(&user).Error

	if err != nil {
		if dbErrors.IsUniqueViolation(err, constraints.UniUsersUserName) {
			c.AbortWithStatusJSON(409, gin.H{"error": "пользователь с таким именем уже существует"})
		} else {
			log.Println(err)
			c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(201, gin.H{"success": "ваш аккаунт успешно создан и ожидает верификации"})

	if _, err := h.NotificationsClient.NotifyAboutNewUserOnVerification(
		c.Request.Context(), &pb.UserOnVerification{ID: uint64(user.ID)},
	); err != nil {
		log.Println(err)
	}
}
