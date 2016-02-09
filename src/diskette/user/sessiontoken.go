package user

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type SessionToken struct {
	Id        string    `json:"id"`
	UserId    string    `json:"userId"`
	CreatedAt time.Time `json:"createdAt"`
}

func (self SessionToken) toString(jwtKey []byte) (string, error) {
	jwtToken := jwt.New(jwt.SigningMethodHS256)
	jwtToken.Claims["id"] = self.Id
	jwtToken.Claims["userId"] = self.UserId
	jwtToken.Claims["createdAt"] = self.CreatedAt.Unix()
	return jwtToken.SignedString(jwtKey)
}

func parseSessionToken(jwtKey []byte, tokenStr string) (token SessionToken, err error) {
	jwtToken, err := jwt.Parse(tokenStr, func(jwtToken *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return
	}
	if !jwtToken.Valid {
		return
	}

	token.Id = jwtToken.Claims["id"].(string)
	token.UserId = jwtToken.Claims["userId"].(string)
	token.CreatedAt = time.Unix(int64(jwtToken.Claims["createdAt"].(float64)), 0)
	return
}
