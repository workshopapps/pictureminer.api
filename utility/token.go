package utility

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

func CreateToken(key, val, secretkey string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["authorized"] = true
	claims["exp"] = time.Now().Add(1 * time.Hour).Unix()
	claims[key] = val

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(secretkey))
	// if err != nil {
	// 	return "", err
	// }

	// return tokenString, nil
}

func DecodeToken(tokenStr, secretkey string) (map[string]interface{}, error) {
	// verify token
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secretkey), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, errors.New("could not verify token")
	}

	return claims, nil
}
