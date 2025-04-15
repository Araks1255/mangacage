package testhelpers

import (
	"time"

	"github.com/Araks1255/mangacage/pkg/common/models"
	"github.com/dgrijalva/jwt-go"
)

func GenerateTokenString(userID uint, secretKey string) (string, error) {
	claims := models.Claims{
		ID: userID,
		StandardClaims: jwt.StandardClaims{
			Subject:   "test",
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
