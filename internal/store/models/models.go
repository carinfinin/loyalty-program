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
	Accrual int64     `json:"accrual,omitempty"`
	User    int       `json:"-"`
	Created time.Time `json:"uploaded_at"`
}

type Balance struct {
	Current   float64 `json:"current"`
	Withdrawn int     `json:"withdrawn"`
}

type Withdrawal struct {
	OrderNumber string    `json:"order"`
	User        int       `json:"-"`
	Sum         int       `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
