package service

import (
	"context"
	"fmt"
	"github.com/carinfinin/loyalty-program/config"
	"github.com/carinfinin/loyalty-program/internal/jwtc"
	"github.com/carinfinin/loyalty-program/internal/logger"
	"github.com/carinfinin/loyalty-program/internal/store"
	"github.com/carinfinin/loyalty-program/internal/store/models"
	"github.com/carinfinin/loyalty-program/internal/store/storepg"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"time"
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

func (s *UserService) Register(ctx context.Context, user *models.User) (string, error) {
	const nf = "service register"

	validate := validator.New()
	err := validate.Struct(user)

	if err != nil {
		logger.Log.Error(nf, fmt.Sprintf("validate error: %v", err))
		return "", err
	}

	passHash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Log.Error(nf, fmt.Sprintf("generate password error: %v", err))
		return "", err
	}

	id, err := s.store.SaveUser(ctx, user.Login, passHash)
	if err != nil {
		logger.Log.Error(nf, fmt.Sprintf(" error: %v", err))
		return "", err
	}

	token, err := jwtc.Generate(user, 24*time.Hour)
	if err != nil {
		return "", err
	}

	logger.Log.Info(nf, fmt.Sprintf(" id: %d", id))
	logger.Log.Info(nf, fmt.Sprintf(" token: %v", token))
	return token, nil

	//todo return cookie
	//todo validate
}
