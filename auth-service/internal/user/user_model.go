package user

import "time"

type User struct {
	ID           string    `json:"ID"`
	Email        string    `json:"Email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"CreatedAt"`
}

func NewUser(email, passwordHash string) *User {
	return &User{Email: email, PasswordHash: passwordHash, CreatedAt: time.Now().UTC()}
}
