package store

import (
	"context"
	"errors"
	"github.com/carinfinin/loyalty-program/internal/store/models"
)

var ErrDouble = errors.New("login already taken")
var ErrNotAuth = errors.New("Invalid login password pair")
var Double = errors.New("rows double")
var Busy = errors.New("uploaded by another user")

type Repository interface {
	User(ctx context.Context, login string) (*models.User, error)
	SaveUser(ctx context.Context, login string, passHash []byte) (int64, error)
	SaveOrder(ctx context.Context, number int64, userID int64) error
	OrderList(ctx context.Context) ([]*models.Order, error)
	WithdrawalSave(ctx context.Context, wd *models.Withdrawal) error
	Withdrawal(ctx context.Context) ([]*models.Withdrawal, error)
}

type OrderRepository interface {
	Order(ctx context.Context, login string, passHash string) (int64, error)
	List(ctx context.Context, login string, passHash string) (int64, error)
}
