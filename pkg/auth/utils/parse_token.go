package utils

import (
	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/dgrijalva/jwt-go"
)

func ParseToken(tokenString, secretKey string) (claims *auth.Claims, err error) {
	token, err := jwt.ParseWithClaims(tokenString, &auth.Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), err
	})

	claims, ok := token.Claims.(*auth.Claims)
	if !ok {
		return nil, err
	}

	return claims, nil
}
