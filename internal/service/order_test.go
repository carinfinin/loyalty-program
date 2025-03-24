package service

import (
	"context"
	"github.com/carinfinin/loyalty-program/internal/config"
	"github.com/carinfinin/loyalty-program/internal/store"
	"github.com/carinfinin/loyalty-program/internal/store/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestSaveOrder(t *testing.T) {

	cfg := config.MewForTest()

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
				repo.On(
					"SaveOrder", mock.Anything, int64(1208), int64(1),
				).Return(int64(1), nil)
				return repo
			}(),
			expectError: false,
		},
		{
			Name: "error",
			Data: data{0, "1406"},
			Repo: func() *mocks.Repository {
				repo := &mocks.Repository{}
				repo.On(
					"SaveOrder", mock.Anything, int64(1406), int64(0),
				).Return(int64(0), store.ErrUserNotFound)
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
