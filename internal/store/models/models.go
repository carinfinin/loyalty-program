package models

import "time"

type User struct {
	ID       int64  `json:"-"`
	Login    string `json:"login" validate:"required,min=3,max=10"`
	Password string `json:"password" validate:"required,min=4,max=10"`
}

type Order struct {
	ID      int64     `json:"-"`
	Number  string    `json:"number"`
	Status  string    `json:"status"`
	Accrual int       `json:"accrual"`
	User    int       `json:"-"`
	Created time.Time `json:"uploaded_at"`
}
