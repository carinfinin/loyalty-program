package storepg

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/carinfinin/loyalty-program/internal/config"
	"github.com/carinfinin/loyalty-program/internal/logger"
	"github.com/carinfinin/loyalty-program/internal/router"
	"github.com/carinfinin/loyalty-program/internal/store"
	"github.com/carinfinin/loyalty-program/internal/store/models"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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
		var errPG *pgconn.PgError
		if errors.As(err, &errPG) && pgerrcode.IsIntegrityConstraintViolation(errPG.Code) {
			return 0, store.ErrDouble
		}

		return 0, err
	}

	return r.RowsAffected()
}

func (s *UserStore) SaveOrder(ctx context.Context, number int64, userID int64) error {
	logger.Log.Debug("start Begin")
	logger.Log.Debug("userID:", userID)

	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var id int64
	row := tx.QueryRowContext(ctx, "SELECT user_id FROM orders WHERE number = $1", number)

	err = row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Debug("no data row start add row")

			_, err = tx.ExecContext(ctx, "INSERT INTO orders (number, status, user_id) VALUES ($1, 'NEW', $2)", number, userID)
			if err != nil {
				logger.Log.Debug("add arder error: ", err)
				return err
			}
			logger.Log.Debug("transaction committed successfully")
			return tx.Commit()
		}
		return err
	}
	if err = row.Err(); err != nil {
		fmt.Println("row err:", err)
		return err
	}
	if id > 0 && id == userID {
		logger.Log.Debug("error add order row double")
		return store.Double
	}
	logger.Log.Debug("error add order row busy")
	return store.Busy

}

func (s *UserStore) OrderList(ctx context.Context) ([]*models.Order, error) {
	const nf = "store order list"
	result := make([]*models.Order, 0)
	userID := ctx.Value(router.UserId)

	rows, err := s.db.QueryContext(ctx, "SELECT number, status, accrual, created_at FROM orders WHERE user_id = $1 ORDER BY created_at DESC", userID)
	if err != nil {
		logger.Log.Debug(nf, fmt.Sprintf("query error: %v", err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		var accrual sql.NullInt64

		err = rows.Scan(&order.Number, &order.Status, &accrual, &order.Created)
		if err != nil {
			logger.Log.Debug(nf, fmt.Sprintf("scan error: %v", err))
			return nil, err
		}
		if accrual.Valid {
			order.Accrual = accrual.Int64
		} else {
			order.Accrual = 0
		}

		result = append(result, &order)
	}
	return result, nil
}
