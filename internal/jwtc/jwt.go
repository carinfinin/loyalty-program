package jwtc

import (
	"fmt"
	"github.com/carinfinin/loyalty-program/internal/config"
	"github.com/carinfinin/loyalty-program/internal/store/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"time"
)

const AuthCookie = "token"

var BadToken = errors.New("bad token")

func Generate(user *models.User, duration time.Duration, cfg *config.Config) (string, error) {
	payload := jwt.MapClaims{
		"uid": user.ID,
		"ulg": user.Login,
		"exp": time.Now().Add(duration).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString([]byte(cfg.Secret))
}

func Decode(token string, cfg *config.Config) (int64, error) {
	claims := jwt.MapClaims{}
	tp, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Неожиданный метод подписи: %v", token.Header["alg"])
		}
		return []byte(cfg.Secret), nil
	})
	if err != nil {
		return 0, err
	}
	if !tp.Valid {
		return 0, fmt.Errorf("is not valid")
	}
	uid, ok := claims["uid"].(float64)
	if !ok {
		return 0, fmt.Errorf("not converce in int64 from %v", claims)
	}
	//jwt.RegisteredClaims
	return int64(uid), nil
}
