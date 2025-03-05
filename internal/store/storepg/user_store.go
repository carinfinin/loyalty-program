package storepg

import (
	"context"
	"github.com/carinfinin/loyalty-program/config"
	"github.com/carinfinin/loyalty-program/internal/logger"
	"github.com/carinfinin/loyalty-program/internal/store/models"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type UserStore struct {
	db *sqlx.DB
}

func New(cfg *config.Config) (*UserStore, error) {
	db, err := sqlx.Open("pgx", cfg.DBPath)
	if err != nil {
		logger.Log.Errorf("store error: %v", err)
		return nil, err
	}

	if err := db.Ping(); err != nil {
		logger.Log.Errorf("failed to ping database: %v", err)
		db.Close()
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
	row.Scan(&user.ID, &user.Password)
	if err := row.Err(); err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserStore) SaveUser(ctx context.Context, login string, passHash []byte) (int64, error) {
	r, err := s.db.ExecContext(ctx, "INSERT INTO users (login, password_hash) VALUES ($1, $2)", login, passHash)
	if err != nil {
		return 0, err
	}
	return r.RowsAffected()
}
