package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/carinfinin/loyalty-program/internal/config"
	"github.com/carinfinin/loyalty-program/internal/logger"
	"github.com/carinfinin/loyalty-program/internal/store"
	"github.com/carinfinin/loyalty-program/internal/store/models"
	"net/http"
	"strconv"
	"time"
)

type ServiceInterface interface {
	Close() error
	Register(ctx context.Context, user *models.User) (string, error)
	Login(ctx context.Context, user *models.User) (string, error)
	Withdrawal(ctx context.Context) ([]*models.Withdrawal, error)
	WithdrawalSave(ctx context.Context, withdrawal *models.Withdrawal) error
	Balance(ctx context.Context) (*models.Balance, error)
	OrderList(ctx context.Context) ([]*models.Order, error)
	SaveOrder(ctx context.Context, number string, userID int64) error
}

type Service struct {
	store      store.Repository
	Config     *config.Config
	chJob      chan *models.Order
	chResult   chan *models.Order
	chBreak    chan struct{}
	retryAfter time.Duration
}

func New(cfg *config.Config, store store.Repository) *Service {
	service := Service{
		store:    store,
		Config:   cfg,
		chJob:    make(chan *models.Order, 2),
		chResult: make(chan *models.Order, 2),
		chBreak:  make(chan struct{}),
	}

	go service.Worker(context.Background())
	go service.Inspector(context.Background())

	return &service
}

func (s *Service) Close() error {

	s.chJob = nil
	s.chResult = nil
	s.chBreak = nil
	return s.store.Close()
}

func (s *Service) Worker(ctx context.Context) {
	semaphore := make(chan struct{}, 10)

	for {
		select {
		case order := <-s.chJob:
			semaphore <- struct{}{}
			go func(o *models.Order) {
				defer func() {
					<-semaphore
				}()
				s.job(o)
			}(order)
		case <-s.chBreak:
			logger.Log.Debug("service Worker s.retryAfter: ", s.retryAfter)
			time.Sleep(s.retryAfter)
		case <-ctx.Done():
			return
		}
	}
}

func (s *Service) job(order *models.Order) {
	const nf = "service job "

	url := fmt.Sprintf("%s/api/orders/%v", s.Config.AccrualAddr, order.Number)
	logger.Log.Info(nf, fmt.Sprintf("url: %v", url))

	response, err := http.Get(url)
	if err != nil {
		logger.Log.Debug(nf, fmt.Sprintf("error: %v", err))
		s.chResult <- order
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusTooManyRequests {
		retryAfter := response.Header.Get("Retry-After")
		logger.Log.Debug(nf, fmt.Sprintf("retry after : %v", retryAfter))
		ra, er := strconv.Atoi(retryAfter)
		if er != nil {
			logger.Log.Debug(nf, fmt.Sprintf("error: %v", err))
		}

		s.retryAfter = time.Duration(ra)
		s.chBreak <- struct{}{}

		s.chResult <- order
		return
	}
	if response.StatusCode == http.StatusNoContent {
		logger.Log.Debug(nf, fmt.Sprintf("order no content status : %v", response.StatusCode))
		time.Sleep(1 * time.Second)
		s.chResult <- order
		return
	}

	var newOrder models.OrderAccrual
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&newOrder)
	if err != nil {
		logger.Log.Debug(nf, fmt.Sprintf("status code : %v", response.StatusCode))
		logger.Log.Error(nf, fmt.Sprintf("order decode error: %v", err))
		s.chResult <- order
		return
	}
	logger.Log.Debug(nf, fmt.Sprintf("new order: %v", newOrder))

	order.Status = newOrder.Status
	order.Accrual = newOrder.Accrual

	s.chResult <- order
}

func (s *Service) Inspector(ctx context.Context) {

	orders := make([]*models.Order, 0, 200)
	ticker := time.NewTicker(5 * time.Second)

	for {
		select {
		case order := <-s.chResult:

			fmt.Println("Inspector order : ", order)

			if order.Status == "PROCESSING" || order.Status == "REGISTERED" || order.Status == "NEW" {
				s.chJob <- order
			}

			if order.Status != "REGISTERED" {
				orders = append(orders, order)
				if len(orders) == 100 {
					err := s.store.OrderBalanceUpdate(context.Background(), orders)
					if err != nil {
						logger.Log.Error("OrderBalanceUpdate error: ", err)
					}
					ticker.Reset(5 * time.Second)
				}
			}
		case <-ctx.Done():
			return
		case <-ticker.C:
			err := s.store.OrderBalanceUpdate(context.Background(), orders)
			if err != nil {
				logger.Log.Error("OrderBalanceUpdate error: ", err)
			}
		default:
			time.Sleep(1 * time.Second)
		}
	}
}
