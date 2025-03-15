package service

import (
	"github.com/carinfinin/loyalty-program/internal/config"
	"github.com/carinfinin/loyalty-program/internal/store"
)

type Service struct {
	store  store.Repository
	Config *config.Config
}

func New(cfg *config.Config, store store.Repository) *Service {
	return &Service{
		store:  store,
		Config: cfg,
	}
}
