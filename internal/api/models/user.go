package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        uint64    `json:"id"`
	Login     string    `json:"login"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
}

func (u *User) HashPwd() error {
	const hashCost = 10
	pass, err := bcrypt.GenerateFromPassword([]byte(u.Password), hashCost)
	if err != nil {
		return err
	}

	u.Password = string(pass)
	return nil
}
