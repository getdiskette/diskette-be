package tokens

import "github.com/dgrijalva/jwt-go"

type ResetToken struct {
	Key string
}

func (self ResetToken) ToString(jwtKey []byte) (string, error) {
	jwtToken := jwt.New(jwt.SigningMethodHS256)
	jwtToken.Claims["key"] = self.Key
	return jwtToken.SignedString(jwtKey)
}

func ParseResetToken(jwtKey []byte, tokenStr string) (token ResetToken, err error) {
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
