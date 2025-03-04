package service

import (
	"context"
	"fmt"
	"github.com/carinfinin/loyalty-program/config"
	"github.com/carinfinin/loyalty-program/internal/logger"
	"github.com/carinfinin/loyalty-program/internal/store"
	"github.com/carinfinin/loyalty-program/internal/store/models"
	"github.com/carinfinin/loyalty-program/internal/store/storepg"
)

type UserService struct {
	store store.UserRepository
}

func New(cfg *config.Config) *UserService {
	s, _ := storepg.New(cfg)
	return &UserService{
		store: s,
	}
}

func (s *UserService) Register(ctx context.Context, user *models.User) error {
	const nf = "service register"
	id, err := s.store.SaveUser(ctx, user.Login, user.PasswordHash)
	if err != nil {
		logger.Log.Error(nf, fmt.Sprintf(" error: %v", err))
		return err
	}

	logger.Log.Info(nf, fmt.Sprintf(" id: %d", id))
	return nil

	//todo return cookie
	//todo validate
}
