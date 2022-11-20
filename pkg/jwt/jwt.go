package jwt

import (
	"errors"

	"github.com/golang-jwt/jwt"
)

// Encode encodes a jwt token using data gotten from payload.
func Encode(secretKey []byte, payload map[string]interface{}) (tokenString string, err error) {
	if len(secretKey) < 1 {
		return "", errors.New("Secret key must be provided !")
	}
	claims := jwt.MapClaims{}
	for k, v := range payload {
		claims[k] = v
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(secretKey)
	return
}

// Decode decodes a jwt token string.
//
// If the jwt token is invalid it returns an error.
func Decode(secretKey []byte, tokenString string) (claims map[string]interface{}, err error) {
	if len(secretKey) < 1 {
		return nil, errors.New("Secret key must be provided !")
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return "", errors.New("Invalid jwt token string")
		}
		return secretKey, nil
	})
	if err != nil {
		return
	}
	if token.Valid {
		claims = token.Claims.(jwt.MapClaims)
		return
	}
	err = errors.New("An unknowm error occured while decoding jwt")
	return
}
