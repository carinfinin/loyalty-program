package router

import (
	"bytes"
	"encoding/json"
	"github.com/carinfinin/loyalty-program/internal/config"
	"github.com/carinfinin/loyalty-program/internal/service/mocks"
	"github.com/carinfinin/loyalty-program/internal/store/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouter_Login(t *testing.T) {
	tests := []struct {
		name        string
		Service     *mocks.ServiceInterface
		cfg         *config.Config
		user        *models.User
		expectError bool
	}{
		{
			"ok",
			func() *mocks.ServiceInterface {
				mockService := &mocks.ServiceInterface{}
				mockService.On("Login", mock.Anything, &models.User{Login: "adm", Password: "err"}).Return("token", nil)
				return mockService
			}(),
			&config.Config{},
			&models.User{Login: "adm", Password: "err"},
			false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			router := Configure(test.cfg, test.Service)

			json, err := json.Marshal(test.user)
			assert.NoError(t, err)
			buf := bytes.NewBuffer(json)

			request := httptest.NewRequest(http.MethodPost, "/api/user/login", buf)
			writer := httptest.NewRecorder()

			router.Login(writer, request)
			result := writer.Result()

			assert.Equal(t, result.StatusCode, 200)
		})
	}

}
