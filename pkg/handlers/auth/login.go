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

	var requestBody struct {
		UserName string `json:"userName" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	var (
		userID   uint
		password string
	)

	row := h.DB.Raw("SELECT id, password FROM users WHERE user_name = ?", requestBody.UserName).Row()

	if err := row.Scan(&userID, &password); err != nil {
		log.Println(err)
	}

	if userID == 0 {
		c.AbortWithStatusJSON(404, gin.H{"error": "аккаунт не найден. возможно, он ещё не прошел верификацию"}) // Сомнительная тема, но я так думаю, что неверифицированный аккаунт итак никаких привелегий не даёт, так что без разницы, войдет в него юзер или нет. Если в будующем это будет не так - поменяю
		return
	}

	if ok := utils.CompareHashPassword(requestBody.Password, password); !ok {
		c.AbortWithStatusJSON(401, gin.H{"error": "неверный пароль"})
		return
	}

	expirationTime := time.Now().Add(2016 * time.Hour)

	claims := models.Claims{
		ID: userID,
		StandardClaims: jwt.StandardClaims{
			Subject:   requestBody.UserName, // Если до сюда дошло, то юзернейм из запроса валидный
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
