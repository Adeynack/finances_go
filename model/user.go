package model

import (
	"crypto/hmac"
	"crypto/md5"
	"encoding/hex"
	"fmt"
)

type User struct {
	BaseModel
	Email             string  `gorm:"not null;uniqueIndex" json:"email"`
	DisplayName       string  `gorm:"not null" json:"display_name"`
	EncryptedPassword string  `gorm:"not null" json:"-"`
	Books             []*Book `gorm:"foreignKey:owner_id" json:"books,omitempty"`
}

func (user *User) SetPassword(password, secret string) error {
	saltedPassword, err := saltValue(secret, password)
	if err != nil {
		return fmt.Errorf("error salting user's password: %w", err)
	}
	user.EncryptedPassword = string(saltedPassword)
	return nil
}

func (user *User) CheckPassword(secret, password string) (bool, error) {
	salted, err := saltValue(secret, password)
	if err != nil {
		return false, fmt.Errorf("error salting password to check: %w", err)
	}
	return hmac.Equal(salted, []byte(user.EncryptedPassword)), nil
}

func saltValue(secret, value string) ([]byte, error) {
	mac := hmac.New(md5.New, []byte(secret))
	_, err := mac.Write([]byte(value))
	if err != nil {
		return nil, fmt.Errorf("error salting value: %w", err)
	}
	return []byte(hex.EncodeToString(mac.Sum(nil))), nil
}
