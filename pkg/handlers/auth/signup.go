package auth

import (
	"context"
	"database/sql"
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/common/utils"
	dbUtils "github.com/Araks1255/mangacage/pkg/common/db/utils"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gin-gonic/gin"
)

func (h handler) Signup(c *gin.Context) {
	_, err := c.Cookie("mangacage_token")
	if err == nil {
		c.AbortWithStatusJSON(409, gin.H{"error": "вы уже вошли в аккаунт"})
		return
	}

	var requestBody struct {
		UserName      string `json:"userName" binding:"required"`
		Password      string `json:"password" binding:"required,min=8"`
		AboutYourself string `json:"aboutYourself"`
	}

	if err = c.ShouldBindJSON(&requestBody); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	user := models.UserOnModeration{
		UserName:      sql.NullString{String: requestBody.UserName, Valid: true},
		AboutYourself: requestBody.AboutYourself,
		Roles:         pq.StringArray([]string{"user"}),
	}

	errChan := make(chan error)

	go func() {
		var errHash error
		user.Password, errHash = utils.GenerateHashPassword(requestBody.Password)
		if errHash != nil {
			log.Println(errHash)
			errChan <- errHash
			return
		}
		errChan <- nil
	}()

	tx := h.DB.Begin()
	defer dbUtils.RollbackOnPanic(tx)
	defer tx.Rollback()

	var existingUserID uint
	tx.Raw("SELECT id FROM users WHERE user_name = ?", user.UserName.String).Scan(&existingUserID)
	if existingUserID != 0 {
		c.AbortWithStatusJSON(403, gin.H{"error": "пользователь с таким именем уже существует"})
		return
	}

	tx.Raw("SELECT id FROM users_on_moderation WHERE user_name = ?", user.UserName.String).Scan(&existingUserID)
	if existingUserID != 0 {
		c.AbortWithStatusJSON(403, gin.H{"error": "пользователь с таким именем уже ожидает верификации"})
		return
	}

	if err = <-errChan; err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if result := tx.Create(&user); result.Error != nil {
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "Ваш аккаунт успешно создан и ожидает верификации"})

	conn, err := grpc.NewClient("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	client := pb.NewNotificationsClient(conn)

	if _, err := client.NotifyAboutUserOnModeration(context.TODO(), &pb.User{ID: uint64(user.ID), New: true}); err != nil {
		log.Println(err)
	}
}
