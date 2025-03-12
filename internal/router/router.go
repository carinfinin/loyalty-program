package router

import (
	"github.com/carinfinin/loyalty-program/internal/config"
	"github.com/carinfinin/loyalty-program/internal/service"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Router struct {
	Handler *chi.Mux
	service *service.Service
	Config  *config.Config
}

func New(cfg *config.Config, service *service.Service) *Router {
	return &Router{
		Handler: chi.NewRouter(),
		service: service,
		Config:  cfg,
	}
}

func Configure(cfg *config.Config, service *service.Service) *Router {
	r := New(cfg, service)
	//r.Use(middleware.RequestID)

	/*
		POST /api/user/orders — загрузка пользователем номера заказа для расчёта;
		GET /api/user/orders — получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
		GET /api/user/balance — получение текущего баланса счёта баллов лояльности пользователя;
		POST /api/user/balance/withdraw — запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
		GET /api/user/withdrawals — получение информации о выводе средств с накопительного счёта пользователем.
	*/
	//r.Handler.Use(middleware.Compress())

	r.Handler.Route("/api/user", func(cr chi.Router) {
		cr.Post("/register", r.RegisterHandler)
		cr.Post("/login", r.LoginHandler)

		cr.With(r.AuthMiddleware).Post("/orders", r.OrderHandler)
		cr.With(r.AuthMiddleware).Get("/orders", r.OrderList)
		cr.With(r.AuthMiddleware).Post("/balance/withdraw", r.Test)
		cr.With(r.AuthMiddleware).Get("/balance", r.Test)
		cr.With(r.AuthMiddleware).Get("/withdrawals", r.Test)
	})
	return r
}

func (r *Router) Test(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("test"))
}
