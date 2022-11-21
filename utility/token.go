package utility

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func CreateToken(key, val, secretkey string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["exp"] = time.Now().Add(1 * time.Hour)
	claims[key] = val

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(secretkey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
