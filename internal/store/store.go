package store

import (
	"context"
	"github.com/carinfinin/loyalty-program/internal/store/models"
)

type Repository interface {
	UserRepository
	OrderRepository
}

type UserRepository interface {
	User(ctx context.Context, login string) (*models.User, error)
	SaveUser(ctx context.Context, login string, passHash string) (int64, error)
}

type OrderRepository interface {
	Order(ctx context.Context, login string, passHash string) (int64, error)
	List(ctx context.Context, login string, passHash string) (int64, error)
}
