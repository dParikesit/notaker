package utils

import (
	"github.com/golang-jwt/jwt"
)

type GoogleMetadata struct {
	Email   string
	Picture string
}

func ExtractGoogle(credential string) (*GoogleMetadata, bool) {

	googleToken, _ := jwt.Parse(credential, nil)

	googleClaims, ok := googleToken.Claims.(jwt.MapClaims)

	if ok {
		return &GoogleMetadata{
			Email: googleClaims["email"].(string),
			Picture: googleClaims["picture"].(string),
		}, true
	}

	return nil, false
}
