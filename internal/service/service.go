package service

import (
	"context"
	"fmt"
	"github.com/carinfinin/loyalty-program/internal/config"
	"github.com/carinfinin/loyalty-program/internal/jwtc"
	"github.com/carinfinin/loyalty-program/internal/logger"
	"github.com/carinfinin/loyalty-program/internal/store"
	"github.com/carinfinin/loyalty-program/internal/store/models"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"strconv"
	"time"
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

	logger.Log.Info(nf, fmt.Sprintf(" id: %d", id))
	logger.Log.Info(nf, fmt.Sprintf(" token: %v", token))
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

func (s *Service) SaveOrder(ctx context.Context, number string) (string, error) {
	const nf = "service save order"
	//id, ok := ctx.Value(router.UserId).(int)
	//if !ok {
	//	return "", fmt.Errorf("user id not get in context")
	//}
	num, err := strconv.ParseInt(number, 10, 64)
	if err != nil {
		logger.Log.Error(nf, err)
		return "", err
	}
	logger.Log.Info("Service SaveOrder")

	_, err = s.store.SaveOrder(ctx, num, 14)
	if err != nil {
		logger.Log.Error(nf, err)
		return "", err
	}
	return "", nil
}
