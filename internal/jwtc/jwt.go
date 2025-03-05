package jwtc

import (
	"github.com/carinfinin/loyalty-program/internal/store/models"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

func Generate(user *models.User, duration time.Duration) (string, error) {
	payload := jwt.MapClaims{
		"uid": user.ID,
		"ulg": user.Login,
		"exp": time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte("hmacSampleSecret"))
}
