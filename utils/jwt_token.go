package utils

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/kistars/pledge-backend/config"
)

func CreateToken(username string) (string, error) {
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24 * 30).Unix(),
	})
	token, err := at.SignedString([]byte(config.Config.Jwt.SecretKey))
	if err != nil {
		return "", err
	}
	return token, nil
}

func ParseToken(token string, secret string) (string, error) {
	claim, err := jwt.Parse(token, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return "", err
	}
	return claim.Claims.(jwt.MapClaims)["username"].(string), nil
}
