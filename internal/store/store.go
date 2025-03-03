package store

import "context"

type UserRepository interface {
	Login(ctx context.Context, login string, passHash string) (int64, error)
	Register(ctx context.Context, login string, passHash string) (int64, error)
}
