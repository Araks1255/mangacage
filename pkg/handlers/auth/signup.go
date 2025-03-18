package auth

import (
	"context"
	"io"
	"log"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/common/utils"

	"github.com/gin-gonic/gin"
)

type UsersProfilePictures struct {
	UserID         uint   `bson:"user_id"`
	ProfilePicture []byte `bson:"profile_picture"`
}

func (h handler) Signup(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	if len(form.Value["userName"]) == 0 || len(form.Value["password"]) == 0 {
		c.AbortWithStatusJSON(400, gin.H{"error": "в запросе недостаточно данных"})
		return
	}

	userName := form.Value["userName"][0]
	password := form.Value["password"][0]

	var aboutYourself string
	if len(form.Value["aboutYourself"]) != 0 {
		aboutYourself = form.Value["aboutYourself"][0]
	}

	user := models.User{
		UserName:      userName,
		AboutYourself: aboutYourself,
	}

	var existingUserID uint
	h.DB.Raw("SELECT id FROM users WHERE user_name = ?", user.UserName).Scan(&existingUserID)
	if existingUserID != 0 {
		c.AbortWithStatusJSON(403, gin.H{"error": "вы уже зарегистрированы"})
		return
	}

	var errHash error
	user.Password, errHash = utils.GenerateHashPassword(password)
	if errHash != nil {
		log.Println(errHash)
		c.AbortWithStatusJSON(500, gin.H{"error": errHash.Error()})
		return
	}

	tx := h.DB.Begin()

	if result := tx.Create(&user); result.Error != nil {
		tx.Rollback()
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	if result := tx.Exec("INSERT INTO user_roles (user_id, role_id) VALUES (?, (SELECT id FROM roles WHERE name = 'user'))", user.ID); result.Error != nil {
		tx.Rollback()
		log.Println(result.Error)
		c.AbortWithStatusJSON(500, gin.H{"error": result.Error.Error()})
		return
	}

	tx.Commit()

	c.JSON(201, gin.H{"success": "регистрация прошла успешно"})

	profilePicture, err := c.FormFile("profilePicture")
	if err != nil {
		log.Println(err)
		return
	}

	file, err := profilePicture.Open()
	if err != nil {
		log.Println(err)
		file.Close()
		return
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		log.Println(err)
		return
	}

	userProfilePicture := UsersProfilePictures{
		UserID:         user.ID,
		ProfilePicture: data,
	}

	if _, err := h.Collection.InsertOne(context.Background(), userProfilePicture); err != nil {
		log.Println(err)
		return
	}
}
