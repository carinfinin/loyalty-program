package service

import (
	"context"
	"fmt"
	"github.com/carinfinin/loyalty-program/internal/config"
	"github.com/carinfinin/loyalty-program/internal/logger"
	"github.com/carinfinin/loyalty-program/internal/store"
	"github.com/carinfinin/loyalty-program/internal/store/models"
	"time"
)

type Service struct {
	store    store.Repository
	Config   *config.Config
	chJob    chan *models.Order
	chResult chan *models.Order
	chBreak  chan struct{}
}

func New(cfg *config.Config, store store.Repository) *Service {
	service := Service{
		store:    store,
		Config:   cfg,
		chJob:    make(chan *models.Order, 100),
		chResult: make(chan *models.Order, 100),
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
			time.Sleep(5 * time.Second)
			fmt.Println(order)

			//request

			order.Status = "PROCESSING"

			s.chResult <- order

		}
	}
	//timer := time.NewTicker(60 * time.Second)
}

func (s *Service) Inspector() {

	for {
		select {
		case order := <-s.chResult:
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
