package service

import (
	"context"
	"github.com/carinfinin/loyalty-program/internal/logger"
	"github.com/carinfinin/loyalty-program/internal/store/models"
	"strconv"
)

func (s *Service) SaveOrder(ctx context.Context, number string, userID int64) error {
	const nf = "service save order"

	num, err := strconv.ParseInt(number, 10, 64)
	if err != nil {
		logger.Log.Error(nf, err)
		return err
	}

	err = s.store.SaveOrder(ctx, num, userID)
	if err != nil {
		logger.Log.Error(nf, err)
		return err
	}

	order := models.Order{
		Number: number,
		User:   userID,
		Status: "NEW",
	}

	s.chJob <- &order

	return nil
}

func (s *Service) OrderList(ctx context.Context) ([]*models.Order, error) {
	return s.store.OrderList(ctx)
}

func (s *Service) Balance(ctx context.Context) (*models.Balance, error) {
	return s.store.Balance(ctx)
}

func (s *Service) WithdrawalSave(ctx context.Context, withdrawal *models.Withdrawal) error {
	return s.store.WithdrawalSave(ctx, withdrawal)
}

func (s *Service) Withdrawal(ctx context.Context) ([]*models.Withdrawal, error) {
	return s.store.Withdrawal(ctx)
}
