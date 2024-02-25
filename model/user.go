package model

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

type User struct {
	BaseModel
	Email             string `gorm:"not null;uniqueIndex" json:"email"`
	DisplayName       string `gorm:"not null" json:"display_name"`
	EncryptedPassword string `gorm:"not null" json:"-"`
}

func (user *User) SetPassword(salt, password string) error {
	saltedPassword, err := saltValue(salt, password)
	if err != nil {
		return fmt.Errorf("error salting user's password: %v", err)
	}
	user.EncryptedPassword = string(saltedPassword)
	return nil
}

func (user *User) CheckPassword(salt, password string) (bool, error) {
	salted, err := saltValue(salt, password)
	if err != nil {
		return false, fmt.Errorf("error salting password to check: %v", err)
	}
	return hmac.Equal(salted, []byte(user.EncryptedPassword)), nil
}

func saltValue(salt, value string) ([]byte, error) {
	mac := hmac.New(md5.New, []byte(salt))
	_, err := mac.Write([]byte(value))
	if err != nil {
		return nil, fmt.Errorf("error salting value: %v", err)
	}
	return []byte(hex.EncodeToString(mac.Sum(nil))), nil
}
