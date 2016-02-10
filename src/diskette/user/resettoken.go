package user

import (
	"github.com/dgrijalva/jwt-go"
)

type resetToken struct {
	key string
}

func (self resetToken) toString(jwtKey []byte) (string, error) {
	jwtToken := jwt.New(jwt.SigningMethodHS256)
	jwtToken.Claims["key"] = self.key
	return jwtToken.SignedString(jwtKey)
}

func parseResetToken(jwtKey []byte, tokenStr string) (token resetToken, err error) {
	jwtToken, err := jwt.Parse(tokenStr, func(jwtToken *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return
	}
	if !jwtToken.Valid {
		return
	}

	token.key = jwtToken.Claims["key"].(string)
	return
}
