package service

import (
	"context"
	"github.com/carinfinin/loyalty-program/internal/config"
	"github.com/carinfinin/loyalty-program/internal/store"
	"github.com/carinfinin/loyalty-program/internal/store/mocks"
	"github.com/carinfinin/loyalty-program/internal/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestService_Register(t *testing.T) {

	cfg := config.NewForTest()

	tests := []struct {
		Name        string
		User        *models.User
		Repo        *mocks.Repository
		Cfg         *config.Config
		expectError bool
	}{
		{
			Name: "OK",
			User: &models.User{
				Login:    "user",
				Password: "234",
			},
			Repo: func() *mocks.Repository {
				repo := &mocks.Repository{}
				repo.On("SaveUser", mock.Anything, "user", mock.Anything).Return(int64(1), nil)
				repo.On("Order", mock.Anything).Return([]*models.Order{}, nil)
				return repo
			}(),
			expectError: false,
		},
		{
			Name: "error",
			User: &models.User{
				Login:    "user",
				Password: "234",
			},
			Repo: func() *mocks.Repository {
				repo := &mocks.Repository{}
				repo.On("SaveUser", mock.Anything, "user", mock.Anything).Return(int64(0), store.ErrDouble)
				repo.On("Order", mock.Anything).Return([]*models.Order{}, nil)
				return repo
			}(),
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			service := New(cfg, test.Repo)
			token, err := service.Register(context.Background(), test.User)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, token)
			}
		})
	}

}

func TestService_Login(t *testing.T) {

	cfg := config.NewForTest()

	tests := []struct {
		Name        string
		User        *models.User
		Repo        *mocks.Repository
		Cfg         *config.Config
		expectError bool
	}{
		{
			Name: "error",
			User: &models.User{
				Login:    "user",
				Password: "234",
			},
			Repo: func() *mocks.Repository {
				repo := &mocks.Repository{}
				repo.On("User", mock.Anything, "user", mock.Anything).Return(nil, store.ErrDouble)
				repo.On("Order", mock.Anything).Return([]*models.Order{}, nil)
				return repo
			}(),
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			service := New(cfg, test.Repo)
			token, err := service.Login(context.Background(), test.User)
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, token)
			}
		})
	}

}
