package utility

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func CreateToken(email string) (string, error) {
	mysecretekey := "Secretkey"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": email,
		"nbf":   time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(mysecretekey))
	if err != nil {
		return "", err
	}

	return tokenString, nil

}
