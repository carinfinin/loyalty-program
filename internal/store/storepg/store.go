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
	var id int64
	err := s.db.QueryRowContext(ctx, "INSERT INTO users (login, password_hash) VALUES ($1, $2) RETURNING id", login, passHash).Scan(&id)
	if err != nil {
		var errPG *pgconn.PgError
		if errors.As(err, &errPG) && pgerrcode.IsIntegrityConstraintViolation(errPG.Code) {
			return 0, store.ErrDouble
		}

		return 0, err
	}

	return id, nil
}

func (s *UserStore) SaveOrder(ctx context.Context, number int64, userID int64) (int64, error) {
	logger.Log.Debug("start Begin")
	logger.Log.Debug("userID:", userID)

	tx, err := s.db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	var id int64
	row := tx.QueryRowContext(ctx, "SELECT user_id FROM orders WHERE number = $1", number)

	err = row.Scan(&id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			logger.Log.Debug("no data row start add row")

			var orderID int64
			err = tx.QueryRowContext(ctx, "INSERT INTO orders (number, status, user_id) VALUES ($1, 'NEW', $2) RETURNING id", number, userID).Scan(&orderID)
			if err != nil {
				logger.Log.Debug("add order error: ", err)
				return 0, err
			}
			err = tx.Commit()
			if err != nil {
				logger.Log.Debug("transaction committed error: ", err)
				return 0, err
			}
			logger.Log.Debug("transaction committed successfully")
			return orderID, nil
		}
		return 0, err
	}
	if err = row.Err(); err != nil {
		fmt.Println("row err:", err)
		return 0, err
	}
	if id > 0 && id == userID {
		logger.Log.Debug("error add order row double")
		return 0, store.ErrRowDouble
	}
	logger.Log.Debug("error add order row busy")
	return 0, store.ErrBusy

}

func (s *UserStore) OrderList(ctx context.Context) ([]*models.Order, error) {
	const nf = "store order list"
	result := make([]*models.Order, 0)
	userID := ctx.Value(router.UserID)

	rows, err := s.db.QueryContext(ctx, "SELECT number, status, accrual, created_at FROM orders WHERE user_id = $1 ORDER BY created_at DESC", userID)
	if err != nil {
		logger.Log.Debug(nf, fmt.Sprintf("query error: %v", err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var order models.Order
		var accrual sql.NullFloat64

		err = rows.Scan(&order.Number, &order.Status, &accrual, &order.Created)
		if err != nil {
			logger.Log.Debug(nf, fmt.Sprintf("scan error: %v", err))
			return nil, err
		}
		if accrual.Valid {
			order.Accrual = accrual.Float64
		} else {
			order.Accrual = 0
		}

		result = append(result, &order)
	}
	err = rows.Err()
	if err != nil {
		logger.Log.Debug(nf, fmt.Sprintf("rows.Err error: %v", err))
		return nil, err
	}
	return result, nil
}

func (s *UserStore) Balance(ctx context.Context) (*models.Balance, error) {
	const nf = "store get balance"
	balance := models.Balance{}
	userID := ctx.Value(router.UserID)

	row := s.db.QueryRowContext(ctx, "SELECT current, withdrawn FROM balance WHERE user_id = $1", userID)
	err := row.Scan(&balance.Current, &balance.Withdrawn)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &balance, nil
		}
		logger.Log.Debug(nf, fmt.Sprintf("scan error: %v", err))
		return nil, err
	}
	err = row.Err()
	if err != nil {
		logger.Log.Debug(nf, fmt.Sprintf("query error: %v", err))
		return nil, err
	}

	return &balance, nil
}

func (s *UserStore) OrderBalanceUpdate(ctx context.Context, order *models.Order) error {
	const nf = "store order and balance update "

	tx, err := s.db.Begin()
	if err != nil {
		logger.Log.Debug(nf, fmt.Sprintf("begin error: %v", err))
		return err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(ctx, "UPDATE orders SET status = $1, accrual = $2 WHERE user_id = $3 AND number = $4", order.Status, order.Accrual, order.User, order.Number)
	if err != nil {
		logger.Log.Debug(nf, fmt.Sprintf("update order error: %v", err))
		return err
	}

	query := `INSERT INTO balance (user_id, current)
				VALUES ($1, $2)
				ON CONFLICT (user_id) DO UPDATE
				SET current = balance.current + EXCLUDED.current;`

	_, err = s.db.ExecContext(ctx, query, order.User, order.Accrual)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (s *UserStore) WithdrawalSave(ctx context.Context, wd *models.Withdrawal) error {
	const nf = "store withdrawal save "

	userID := ctx.Value(router.UserID)

	tx, err := s.db.Begin()
	if err != nil {
		logger.Log.Debug(nf, fmt.Sprintf("begin error: %v", err))
		return err
	}
	defer tx.Rollback()

	//get
	balance := models.Balance{}
	row := tx.QueryRowContext(ctx, "SELECT id, current, withdrawn FROM balance WHERE user_id = $1", userID)
	err = row.Scan(&balance.ID, &balance.Current, &balance.Withdrawn)
	if err != nil {
		logger.Log.Debug(nf, fmt.Sprintf("get balance error: %v", err))
		return err
	}

	//compare
	if balance.Current < wd.Sum {
		return store.ErrBalanceLow
	}
	balance.Current = balance.Current - wd.Sum
	balance.Withdrawn = balance.Withdrawn + wd.Sum

	//update
	_, err = tx.ExecContext(ctx, "UPDATE balance SET current = $1, withdrawn = $2 WHERE id = $3 AND user_id = $4", balance.Current, balance.Withdrawn, balance.ID, userID)
	if err != nil {
		logger.Log.Debug(nf, fmt.Sprintf("update balance error: %v", err))
		return err
	}

	_, err = s.db.ExecContext(ctx, "INSERT INTO withdrawals (order_number, sum, user_id) VALUES ($1, $2, $3)", wd.OrderNumber, wd.Sum, userID)
	if err != nil {
		logger.Log.Debug(nf, fmt.Sprintf("scan error: %v", err))
		return err
	}
	return tx.Commit()
}

func (s *UserStore) Withdrawal(ctx context.Context) ([]*models.Withdrawal, error) {

	const nf = "store get withdrawal "
	userID := ctx.Value(router.UserID)
	result := make([]*models.Withdrawal, 0)
	rows, err := s.db.QueryContext(ctx, "SELECT order_number, sum, processed_at FROM withdrawals WHERE user_id = $1 ORDER BY processed_at DESC", userID)
	if err != nil {
		logger.Log.Debug(nf, fmt.Sprintf("query error: %v", err))
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		tmp := models.Withdrawal{}
		err = rows.Scan(&tmp.OrderNumber, &tmp.Sum, &tmp.ProcessedAt)
		if err != nil {
			logger.Log.Debug(nf, fmt.Sprintf("scan error: %v", err))
			return nil, err
		}
		result = append(result, &tmp)
	}
	err = rows.Err()
	if err != nil {
		logger.Log.Debug(nf, fmt.Sprintf("rows.Err error: %v", err))
		return nil, err
	}
	return result, nil
}
func (s *UserStore) Close() error {
	return s.db.Close()
}
