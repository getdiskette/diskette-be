package tokens

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type SessionToken struct {
	UserId    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}

func (self SessionToken) ToString(jwtKey []byte) (string, error) {
	jwtToken := jwt.New(jwt.SigningMethodHS256)
	jwtToken.Claims["userId"] = self.UserId
	jwtToken.Claims["createdAt"] = self.CreatedAt.Unix()
	return jwtToken.SignedString(jwtKey)
}

func ParseSessionToken(jwtKey []byte, tokenStr string) (token SessionToken, err error) {
	jwtToken, err := jwt.Parse(tokenStr, func(jwtToken *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return
	}
	if !jwtToken.Valid {
		return
	}

	token.UserId = jwtToken.Claims["userId"].(string)
	token.CreatedAt = time.Unix(int64(jwtToken.Claims["createdAt"].(float64)), 0)
	return
}
