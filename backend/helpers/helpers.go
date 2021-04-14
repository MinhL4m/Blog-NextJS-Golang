package helpers

import (
	"errors"

	"github.com/badoux/checkmail"
	"golang.org/x/crypto/bcrypt"
)

/**----------Password Related Functions-------------*/
func Hash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

/**----------Validation Functions-------------*/
func UpdateValid(username, password, email string) error {
	if username == "" {
		return errors.New("required username")
	}
	if password == "" {
		return errors.New("required password")
	}
	if email == "" {
		return errors.New("required email")
	}
	if err := checkmail.ValidateFormat(email); err != nil {
		return errors.New("invalid email")
	}

	return nil
}

func LoginValid(password, email string) error {
	if password == "" {
		return errors.New("required password")
	}
	if email == "" {
		return errors.New("required email")
	}
	if err := checkmail.ValidateFormat(email); err != nil {
		return errors.New("invalid email")
	}

	return nil
}