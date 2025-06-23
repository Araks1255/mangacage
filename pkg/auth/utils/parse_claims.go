package utils

import "github.com/Araks1255/mangacage/pkg/auth"

func ParseClaims(cookieFn func(string) (string, error), secretKey string) (*auth.Claims, error) {
	cookie, err := cookieFn("mangacage_token")
	if err != nil {
		return nil, err
	}

	claims, err := ParseToken(cookie, secretKey)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
