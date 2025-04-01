package router

import (
	"github.com/carinfinin/loyalty-program/internal/config"
	"github.com/carinfinin/loyalty-program/internal/service"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Router struct {
	Handler *chi.Mux
	service service.ServiceInterface
	Config  *config.Config
}

func New(cfg *config.Config, service service.ServiceInterface) *Router {
	return &Router{
		Handler: chi.NewRouter(),
		service: service,
		Config:  cfg,
	}
}

func Configure(cfg *config.Config, service service.ServiceInterface) *Router {
	r := New(cfg, service)

	r.Handler.Use(middleware.Compress(5, "text/html", "text/plain", "application/json"))
	r.Handler.Route("/api/user", func(cr chi.Router) {
		cr.Post("/register", r.Register)
		cr.Post("/login", r.Login)

		cr.With(r.AuthMiddleware).Post("/orders", r.OrderSave)
		cr.With(r.AuthMiddleware).Get("/orders", r.OrderList)
		cr.With(r.AuthMiddleware).Post("/balance/withdraw", r.WithdrawalSave)
		cr.With(r.AuthMiddleware).Get("/balance", r.Balance)
		cr.With(r.AuthMiddleware).Get("/withdrawals", r.Withdrawal)
	})
	return r
}
