package models

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	UserKey = "UserKey"
)

type User struct {
	ID                string `json:"id" bson:"_id" storm:"id,increment"`
	Email             string `json:"email" bson:"email" storm:"unique" validate:"required,email"`
	Name              string `json:"name"`
	Password          string `json:"password,omitempty" validate:"required,min=6,max=100"`
	EncryptedPassword string `json:"encrypted_password,omitempty" bson:"encrypted_password"`
	//Role              string `json:"-"`
	//AclGroup          string `json:"aclGroup,omitempty"`
}

func (u *User) CheckPermission(action Action) bool {
	//if u.Email == "example@mail.com" {
	//	return false
	//}
	return true
}

//// Validate ...
//func (u *User) Validate() error {
//	return validation.ValidateStruct(
//		u,
//		validation.Field(&u.Email, validation.Required, is.Email),
//		validation.Field(&u.Password, validation.By(requiredIf(u.EncryptedPassword == "")), validation.Length(6, 100)),
//	)
//}

// BeforeCreate ...
func (u *User) BeforeCreate() error {
	if len(u.Password) > 0 {
		enc, err := encryptString(u.Password)
		if err != nil {
			return err
		}

		u.EncryptedPassword = enc
	}

	return nil
}

// Sanitize ...
func (u *User) Sanitize() {
	u.Password = ""
}

//func requiredIf(cond bool) validation.RuleFunc {
//	return func(value interface{}) error {
//		if cond {
//			return validation.Validate(value, validation.Required)
//		}
//
//		return nil
//	}
//}

// ComparePassword ...
func (u *User) ComparePassword(password string) bool {
	return bcrypt.CompareHashAndPassword([]byte(u.EncryptedPassword), []byte(password)) == nil
}

func encryptString(s string) (string, error) {
	b, err := bcrypt.GenerateFromPassword([]byte(s), bcrypt.MinCost)
	if err != nil {
		return "", err
	}

	return string(b), nil
}
