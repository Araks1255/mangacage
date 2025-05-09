package auth

import (
	"database/sql"
	"log"
	"time"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/auth/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func (h handler) Login(c *gin.Context) {
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
		userID       sql.NullInt64
		passwordHash sql.NullString
	)

	row := h.DB.Raw("SELECT id, password FROM users WHERE user_name = ?", requestBody.UserName).Row()

	if err := row.Scan(&userID, &passwordHash); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if !userID.Valid {
		c.AbortWithStatusJSON(404, gin.H{"error": "аккаунт не найден. возможно, он ещё не прошел верификацию"}) // Сомнительная тема, но я так думаю, что неверифицированный аккаунт итак никаких привелегий не даёт, так что без разницы, войдет в него юзер или нет. Если в будующем это будет не так - поменяю
		return
	}

	if ok := utils.CompareHashPassword(requestBody.Password, passwordHash.String); !ok {
		c.AbortWithStatusJSON(401, gin.H{"error": "неверный пароль"})
		return
	}

	expirationTime := time.Now().Add(2016 * time.Hour)

	claims := auth.Claims{
		ID: uint(userID.Int64),
		StandardClaims: jwt.StandardClaims{
			Subject:   requestBody.UserName, // Если до сюда дошло, то юзернейм из запроса валидный
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(h.SecretKey))
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("mangacage_token", tokenString, int(expirationTime.Unix()), "/", "localhost", false, true) // ПОМЕНЯТЬ НА ПРОДЕ

	c.JSON(200, gin.H{"success": "Вход в аккаунт выполнен успешно"})
}
