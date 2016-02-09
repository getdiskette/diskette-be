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
	token := jwt.New(jwt.SigningMethodHS256)
	token.Claims["id"] = self.Id
	token.Claims["userId"] = self.UserId
	token.Claims["createdAt"] = self.CreatedAt.Unix()
	return token.SignedString(jwtKey)
}

func parseSessionToken(jwtKey []byte, sessionTokenStr string) (sessionToken SessionToken, err error) {
	token, err := jwt.Parse(sessionTokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return
	}
	if !token.Valid {
		return
	}

	sessionToken.Id = token.Claims["id"].(string)
	sessionToken.UserId = token.Claims["userId"].(string)
	sessionToken.CreatedAt = time.Unix(int64(token.Claims["createdAt"].(float64)), 0)
	return
}
