package models

type User struct {
	ID           int64
	Login        string `json:"login"`
	PasswordHash string `json:"password"`
}
