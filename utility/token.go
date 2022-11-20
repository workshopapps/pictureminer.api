package utility

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func CreateToken(email string) (string, error) {
	mysecretekey := "Secretkey"
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["exp"] = time.Now().Add(1 * time.Hour)
	claims["email"] = email

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(mysecretekey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
