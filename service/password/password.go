package password

import "golang.org/x/crypto/bcrypt"

func PasswordIsValid(userPassword, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(userPassword), []byte(providedPassword))

	check := true
	msg := ""

	if err != nil {
		check = false
		msg = "invalid email or password"
	}
	return check, msg
}
