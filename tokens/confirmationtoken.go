package tokens

import "github.com/dgrijalva/jwt-go"

type ConfirmationToken struct {
	Key string
}

func (service ConfirmationToken) ToString(jwtKey []byte) (string, error) {
	jwtToken := jwt.New(jwt.SigningMethodHS256)
	jwtToken.Claims["key"] = service.Key
	return jwtToken.SignedString(jwtKey)
}

func ParseConfirmationToken(jwtKey []byte, tokenStr string) (token ConfirmationToken, err error) {
	jwtToken, err := jwt.Parse(tokenStr, func(jwtToken *jwt.Token) (interface{}, error) {
		return []byte(jwtKey), nil
	})
	if err != nil {
		return
	}
	if !jwtToken.Valid {
		return
	}

	token.Key = jwtToken.Claims["key"].(string)
	return
}
