package service

import (
	"context"
	"fmt"
	"github.com/carinfinin/loyalty-program/internal/jwtc"
	"github.com/carinfinin/loyalty-program/internal/logger"
	"github.com/carinfinin/loyalty-program/internal/store"
	"github.com/carinfinin/loyalty-program/internal/store/models"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"time"
)

func (s *Service) Register(ctx context.Context, user *models.User) (string, error) {
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
	user.ID = id

	token, err := jwtc.Generate(user, 24*time.Hour, s.Config)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) Login(ctx context.Context, user *models.User) (string, error) {
	const nf = "service login"

	validate := validator.New()
	err := validate.Struct(user)

	if err != nil {
		logger.Log.Error(nf, fmt.Sprintf("validate error: %v", err))
		return "", err
	}

	u, err := s.store.User(ctx, user.Login)
	if err != nil {
		logger.Log.Error(nf, fmt.Sprintf(" error: %v", err))
		return "", store.ErrNotAuth
	}

	if err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(user.Password)); err != nil {
		return "", store.ErrNotAuth
	}

	token, err := jwtc.Generate(u, 24*time.Hour, s.Config)
	if err != nil {
		return "", err
	}

	return token, nil
}
