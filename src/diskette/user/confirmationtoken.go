package user

import (
	"github.com/dgrijalva/jwt-go"
)

type privateConfirmationToken struct {
	database string
	key      string
	language string
}

func (self privateConfirmationToken) toString(jwtKey []byte) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["database"] = self.database
	token.Claims["key"] = self.key
	token.Claims["language"] = self.language
	return token.SignedString(jwtKey)
}

func parseConfirmationToken(jwtKey []byte, confirmationTokenStr string) (confirmationToken privateConfirmationToken, err error) {
	token, err := jwt.Parse(confirmationTokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return
	}
	if !token.Valid {
		return
	}

	confirmationToken.database = token.Claims["database"].(string)
	confirmationToken.key = token.Claims["key"].(string)
	confirmationToken.language = token.Claims["language"].(string)
	return
}
