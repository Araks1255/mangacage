package auth

import (
	"log"
	"time"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/Araks1255/mangacage/pkg/auth/utils"
	"github.com/Araks1255/mangacage/pkg/common/models/dto"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func (h handler) Login(c *gin.Context) {
	if _, err := c.Cookie("mangacage_token"); err == nil {
		c.AbortWithStatusJSON(401, gin.H{"error": "вы уже вошли в аккаунт"})
		return
	}

	var requestBody dto.CreateUserDTO

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(400, gin.H{"error": err.Error()})
		return
	}

	var check struct {
		UserID       *uint
		PasswordHash string
	}

	err := h.DB.Raw(
		"SELECT id AS user_id, password AS password_hash FROM users WHERE user_name = ?",
		requestBody.UserName,
	).Scan(&check).Error

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	if check.UserID == nil {
		c.AbortWithStatusJSON(401, gin.H{"error": "аккаунт не найден"})
		return
	}

	if !utils.CompareHashPassword(requestBody.Password, check.PasswordHash) {
		c.AbortWithStatusJSON(401, gin.H{"error": "неверный пароль"})
		return
	}

	expires := time.Now().Add(60 * 64 * time.Hour)
	maxAge := int(time.Until(expires).Seconds())

	claims := auth.Claims{
		ID: *check.UserID,
		StandardClaims: jwt.StandardClaims{
			Subject:   requestBody.UserName,
			ExpiresAt: expires.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(h.SecretKey))
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(500, gin.H{"error": err.Error()})
		return
	}

	c.SetCookie("mangacage_token", tokenString, maxAge, "/", h.Host, false, true) // ПОМЕНЯТЬ НА ПРОДЕ

	c.JSON(200, gin.H{"success": "Вход в аккаунт выполнен успешно"})
}
