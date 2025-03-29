package service

import (
	"context"
	"encoding/json"
	"github.com/carinfinin/loyalty-program/internal/config"
	"github.com/carinfinin/loyalty-program/internal/store"
	"github.com/carinfinin/loyalty-program/internal/store/mocks"
	"github.com/carinfinin/loyalty-program/internal/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSaveOrder(t *testing.T) {

	cfg := config.NewForTest()

	type data struct {
		userID int64
		Order  string
	}
	tests := []struct {
		Name        string
		Data        data
		Repo        *mocks.Repository
		Cfg         *config.Config
		expectError bool
	}{
		{
			Name: "OK",
			Data: data{1, "1208"},
			Repo: func() *mocks.Repository {
				repo := &mocks.Repository{}
				repo.On("SaveOrder", mock.Anything, int64(1208), int64(1)).Return(int64(1), nil)
				repo.On("Order", mock.Anything).Return([]*models.Order{}, nil)
				return repo
			}(),
			expectError: false,
		},
		{
			Name: "error",
			Data: data{0, "1406"},
			Repo: func() *mocks.Repository {
				repo := &mocks.Repository{}
				repo.On("SaveOrder", mock.Anything, int64(1406), int64(0)).Return(int64(0), store.ErrUserNotFound)
				repo.On("Order", mock.Anything).Return([]*models.Order{}, nil)
				return repo
			}(),
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {

			service := New(cfg, test.Repo)
			err := service.SaveOrder(context.Background(), test.Data.Order, int64(test.Data.userID))
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}

func TestSaveOrder2(t *testing.T) {

	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		assert.Equal(t, "/api/orders/0109", request.URL.Path)
		assert.Equal(t, "GET", request.Method)

		response := struct {
			Order   string `json:"order"`
			Status  string `json:"status"`
			Accrual int    `json:"accrual"`
		}{
			"2109",
			"PROCESSED",
			500,
		}
		writer.Header().Set("Content-Type", "application/json")
		writer.WriteHeader(http.StatusOK)
		json.NewEncoder(writer).Encode(response)
	}))
	defer testServer.Close()

	cfg := config.NewForTest()

	cfg.AccrualAddr = testServer.URL

	type data struct {
		userID int64
		Order  string
	}
	tests := []struct {
		Name        string
		Data        data
		Repo        *mocks.Repository
		Cfg         *config.Config
		expectError bool
	}{
		{
			Name: "OK",
			Data: data{1, "2109"},
			Repo: func() *mocks.Repository {
				repo := &mocks.Repository{}
				repo.On("SaveOrder", mock.Anything, int64(2109), int64(1)).Return(int64(1), nil)
				repo.On("Order", mock.Anything).Return([]*models.Order{}, nil)
				return repo
			}(),
			expectError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {

			service := New(cfg, test.Repo)
			err := service.SaveOrder(context.Background(), test.Data.Order, int64(test.Data.userID))
			if test.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
