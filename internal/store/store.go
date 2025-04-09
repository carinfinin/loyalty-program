package store

import (
	"context"
	"errors"
	"github.com/carinfinin/loyalty-program/internal/store/models"
)

var ErrDouble = errors.New("login already taken")
var ErrNotAuth = errors.New("invalid login password pair")
var ErrRowDouble = errors.New("rows double")
var ErrBusy = errors.New("uploaded by another user")
var ErrBalanceLow = errors.New("there are insufficient funds in the account")
var ErrUserNotFound = errors.New("user if not found")

type Repository interface {
	User(ctx context.Context, login string) (*models.User, error)
	SaveUser(ctx context.Context, login string, passHash []byte) (int64, error)
	SaveOrder(ctx context.Context, number int64, userID int64) (int64, error)
	OrderList(ctx context.Context) ([]*models.Order, error)
	WithdrawalSave(ctx context.Context, wd *models.Withdrawal) error
	Withdrawal(ctx context.Context) ([]*models.Withdrawal, error)
	Balance(ctx context.Context) (*models.Balance, error)
	OrderBalanceUpdate(ctx context.Context, orders []*models.Order) error
	Order(ctx context.Context) ([]*models.Order, error)
	Close() error
}
