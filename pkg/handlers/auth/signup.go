package auth

import (
	"context"
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/common/utils"
	pb "github.com/Araks1255/mangacage_protos"
	"github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gin-gonic/gin"
)

type UsersProfilePictures struct {
	UserID         uint   `bson:"user_id"`
	ProfilePicture []byte `bson:"profile_picture"`
}

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
		UserName:      requestBody.UserName,
		AboutYourself: requestBody.AboutYourself,
		Roles:         pq.StringArray([]string{"user"}),
	}

	user.Password, err = utils.GenerateHashPassword(requestBody.Password) // Генерация хэша происходит заранее, чтобы не делать этого в транзации, блокируя запись в бд на секунду
	if err != nil {
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	tx := h.DB.Begin()
	if r := recover(); r != nil {
		tx.Rollback()
		panic(r)
	}

	var existingUserID uint
	tx.Raw("SELECT id FROM users WHERE user_name = ?", user.UserName).Scan(&existingUserID)
	if existingUserID != 0 {
		tx.Rollback()
		c.AbortWithStatusJSON(403, gin.H{"error": "пользователь с таким именем уже существует"})
		return
	}

	tx.Raw("SELECT id FROM users_on_moderation WHERE user_name = ?", user.UserName).Scan(&existingUserID)
	if existingUserID != 0 {
		tx.Rollback()
		c.AbortWithStatusJSON(403, gin.H{"error": "пользователь с таким именем уже ожидает верификации"})
		return
	}

	if result := tx.Create(&user); result.Error != nil {
		tx.Rollback()
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

	if _, err := client.NotifyAboutUser(context.Background(), &pb.User{Name: user.UserName}); err != nil {
		log.Println(err)
	}
}
