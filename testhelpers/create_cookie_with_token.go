package testhelpers

import (
	"net/http"

	"time"

	"github.com/Araks1255/mangacage/pkg/auth"
	"github.com/dgrijalva/jwt-go"
)

func CreateCookieWithToken(userID uint, secretKey string) (*http.Cookie, error) {
	claims := auth.Claims{
		ID: userID,
		StandardClaims: jwt.StandardClaims{
			Subject:   "test",
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return nil, err
	}

	cookie := &http.Cookie{
		Name:  "mangacage_token",
		Value: tokenString,
	}

	return cookie, nil
}
