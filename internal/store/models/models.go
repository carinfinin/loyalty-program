package models

type User struct {
	ID       int64  `json:"-"`
	Login    string `json:"login" validate:"required,min=3,max=10"`
	Password string `json:"password" validate:"required,min=4,max=10"`
}
