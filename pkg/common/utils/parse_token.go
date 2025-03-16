package utils

import (
	"github.com/Araks1255/mangacage/pkg/common/models"

	"github.com/dgrijalva/jwt-go"
)

func ParseToken(tokenString, secretKey string) (claims *models.Claims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), err
	})

	claims, ok := token.Claims.(*models.Claims)
	if !ok {
		return nil, err
	}

	return claims, nil
}
