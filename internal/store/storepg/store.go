package storepg

import (
	"context"
	"github.com/carinfinin/loyalty-program/config"
	"github.com/carinfinin/loyalty-program/internal/store/models"
	"github.com/carinfinin/loyalty-program/pkg/logger"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type UserStore struct {
	db *sqlx.DB
}

func New(cfg config.Config) (*UserStore, error) {
	db, err := sqlx.Open("pgx", cfg.DBPath)
	if err != nil {
		logger.Log.Error("store error: ", err)
		return nil, err
	}
	return &UserStore{
		db: db,
	}, nil
}

func (s *UserStore) User(ctx context.Context, login string) (*models.User, error) {
	user := models.User{
		Login: login,
	}
	row := s.db.QueryRowContext(ctx, "SELECT id, password_hash FROM users WHERE login = $1", login)
	row.Scan(&user.ID, &user.PasswordHash)
	if err := row.Err(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserStore) SaveUser(ctx context.Context, login, passHash string) (int64, error) {
	r, err := s.db.ExecContext(ctx, "INSERT INTO users (login, password_hash) VALUES ($1, $2)", login, passHash)
	if err != nil {
		return 0, err
	}
	return r.RowsAffected()
}
