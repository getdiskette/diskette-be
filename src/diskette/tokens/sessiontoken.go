package tokens

import (
	"time"

	"github.com/dgrijalva/jwt-go"
)

type SessionToken struct {
	UserID    string    `json:"UserID"`
	CreatedAt time.Time `json:"createdAt"`
}

func (service SessionToken) ToString(jwtKey []byte) (string, error) {
	jwtToken := jwt.New(jwt.SigningMethodHS256)
	jwtToken.Claims["UserID"] = service.UserID
	jwtToken.Claims["createdAt"] = service.CreatedAt.Unix()
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

	token.UserID = jwtToken.Claims["UserID"].(string)
	token.CreatedAt = time.Unix(int64(jwtToken.Claims["createdAt"].(float64)), 0)
	return
}
