package users

import (
	"golang.org/x/crypto/bcrypt"
)

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (u *User) CheckUsersPasswordHash(passwordHash string) bool {
	return u.PasswordHash == passwordHash
}

func CheckUsersPassword(email, password string) bool {
	u, err := store.getUserByEmail(email)
	if err != nil {
		return false
	}
	return CheckPasswordHash(password, u.PasswordHash)
}
