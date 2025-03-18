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

	go service.Worker()
	go service.Inspector()

	return &service
}

func (s *Service) Close() error {
	return s.store.Close()
}

func (s *Service) Worker() {

	for {
		select {
		case order := <-s.chJob:
			go s.job(order)
		case <-s.chBreak:
			logger.Log.Debug("service Worker s.retryAfter: ", s.retryAfter)
			time.Sleep(s.retryAfter)
		}
	}
}

func (s *Service) job(order *models.Order) {
	const nf = "service job "

	url := fmt.Sprintf("%s/api/orders/%d", s.Config.AccrualAddr, order.ID)
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
		s.chResult <- order
		return
	}

	var newOrder models.Order
	decoder := json.NewDecoder(response.Body)
	err = decoder.Decode(&newOrder)
	if err != nil {
		logger.Log.Debug(nf, fmt.Sprintf("status code : %v", response.StatusCode))
		logger.Log.Error(nf, fmt.Sprintf("order decode error: %v", err))
		s.chResult <- order
		return
	}
	logger.Log.Debug(nf, fmt.Sprintf("new order: %v", newOrder))
	s.chResult <- &newOrder

}

func (s *Service) Inspector() {

	for {
		select {
		case order := <-s.chResult:

			fmt.Println("Inspector order : ", order)

			if order.Status == "PROCESSING" || order.Status == "REGISTERED" {
				s.chJob <- order
			}

			if order.Status != "REGISTERED" {
				err := s.store.OrderBalanceUpdate(context.Background(), order)
				if err != nil {
					logger.Log.Error("OrderBalanceUpdate error: ", err)
				}
			}
		}
	}
}
