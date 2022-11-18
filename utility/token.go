package utility

import (
	"github.com/golang-jwt/jwt/v4"
	"log"
	"time"
)

func CreateToken() *string {
	mysecretekey := "Secretkey"
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"foo": "bar",
		"nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
	})
	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(mysecretekey))
	if err != nil {
		log.Panic("Error is getting jwt key")
	}

	return &tokenString

}
