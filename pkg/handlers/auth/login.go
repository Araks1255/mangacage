package auth

import (
	"log"
	"time"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/Araks1255/mangacage/pkg/common/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func (h handler) Login(c *gin.Context) {
	viper.SetConfigFile("./pkg/common/envs/.env")
	viper.ReadInConfig()

	secretKey := []byte(viper.Get("SECRET_KEY").(string))

	var user models.User

	if err := c.ShouldBindJSON(&user); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	var existingUser models.User
	h.DB.Raw("SELECT * FROM users WHERE user_name = ?", user.UserName).Scan(&existingUser)
	if existingUser.ID == 0 {
		c.AbortWithStatusJSON(401, gin.H{"error": "Аккаунт не найден"})
		return
	}

	if ok := utils.CompareHashPassword(user.Password, existingUser.Password); !ok {
		c.AbortWithStatusJSON(401, gin.H{"error": "Неверный пароль"})
		return
	}

	expirationTime := time.Now().Add(2016 * time.Hour)

	claims := models.Claims{
		ID:   existingUser.ID,
		Role: "user",
		StandardClaims: jwt.StandardClaims{
			Subject:   existingUser.UserName,
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("mangacage_token", tokenString, int(expirationTime.Unix()), "/", "localhost", false, true)

	c.JSON(200, gin.H{"success": "Вход в аккаунт выполнен успешно"})
}
