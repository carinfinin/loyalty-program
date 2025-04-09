package jwtc

import (
	"github.com/carinfinin/loyalty-program/internal/config"
	"github.com/carinfinin/loyalty-program/internal/store/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGenerateDecode(t *testing.T) {
	tests := []struct {
		name        string
		user        *models.User
		cfg         *config.Config
		duration    time.Duration
		expectError bool
	}{
		{
			"OK",
			&models.User{ID: 122},
			&config.Config{Secret: "SECRET"},
			20 * time.Second,
			false,
		},
		{
			"OK",
			&models.User{ID: 1568},
			&config.Config{Secret: "SECRET"},
			20 * time.Second,
			false,
		},
		{
			"Empty secret",
			&models.User{ID: 122},
			&config.Config{Secret: ""},
			0 * time.Second,
			true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			token, err := Generate(test.user, test.duration, test.cfg)
			assert.NoError(t, err)
			assert.NotEmpty(t, token, "токен должен быть не пустой")

			id, err := Decode(token, test.cfg)
			if test.expectError {
				assert.Error(t, err, "Ожидалась ошибка")
				return
			}
			assert.NoError(t, err)
			assert.NotEmpty(t, id, "id должен быть больеше 0")
			assert.Equal(t, id, test.user.ID)
		})
	}
}
