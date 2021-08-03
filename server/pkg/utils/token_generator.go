package utils

import (
	"time"

	"github.com/golang-jwt/jwt"

)

func GenerateAccess (email string, exp time.Time) (string, error) {
	secret := Dotenv("ACCESS")

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["email"] = email
	claims["exp"] = exp

	access, err := token.SignedString([]byte(secret))
	if err!=nil {
		return "", err
	}
	
	return access, nil
}

func GenerateRefresh (email string) (string, error) {
	secret := Dotenv("REFRESH")

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["email"] = email
	claims["exp"] = 0

	refresh, err := token.SignedString([]byte(secret))
	if err!=nil {
		return "", err
	}
	
	return refresh, nil
}