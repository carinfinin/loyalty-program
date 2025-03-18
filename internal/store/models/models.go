package models

import "time"

type User struct {
	ID       int64  `json:"-"`
	Login    string `json:"login" validate:"required,min=1,max=100"`
	Password string `json:"password" validate:"required,min=1,max=100"`
}

type Order struct {
	ID      int64     `json:"-"`
	Number  string    `json:"number"`
	Status  string    `json:"status"`
	Accrual float64   `json:"accrual,omitempty"`
	User    int64     `json:"-"`
	Created time.Time `json:"uploaded_at"`
}

type OrderAccrual struct {
	Number  string  `json:"order"`
	Status  string  `json:"status"`
	Accrual float64 `json:"accrual,omitempty"`
}

type Balance struct {
	ID        int64     `json:"-"`
	Current   float64   `json:"current"`
	Withdrawn float64   `json:"withdrawn"`
	UserID    int64     `json:"-"`
	CreatedAt time.Time `json:"-"`
}

type Withdrawal struct {
	OrderNumber string    `json:"order"`
	User        int64     `json:"-"`
	Sum         float64   `json:"sum"`
	ProcessedAt time.Time `json:"processed_at"`
}
