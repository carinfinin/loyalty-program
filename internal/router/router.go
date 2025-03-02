package router

import (
	"github.com/carinfinin/loyalty-program/config"
	"github.com/go-chi/chi/v5"
	"net/http"
)

type Router struct {
	Handler *chi.Mux
}

func New(cfg *config.Config) *Router {
	return &Router{
		Handler: chi.NewRouter(),
	}
}

func Configure(cfg *config.Config) *Router {
	r := New(cfg)
	//r.Use(middleware.RequestID)
	r.Handler.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	/*
		POST /api/user/orders — загрузка пользователем номера заказа для расчёта;
		GET /api/user/orders — получение списка загруженных пользователем номеров заказов, статусов их обработки и информации о начислениях;
		GET /api/user/balance — получение текущего баланса счёта баллов лояльности пользователя;
		POST /api/user/balance/withdraw — запрос на списание баллов с накопительного счёта в счёт оплаты нового заказа;
		GET /api/user/withdrawals — получение информации о выводе средств с накопительного счёта пользователем.
	*/
	r.Handler.Route("/api/user", func(cr chi.Router) {
		cr.Post("/register", r.RegisterHandler)
		cr.Post("/login", r.Test)
		cr.Post("/orders", r.Test)
		cr.Post("/balance/withdraw", r.Test)
		cr.Get("/balance", r.Test)
		cr.Get("/withdrawals", r.Test)
	})
	return r
}

func (r *Router) Test(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("test"))
}
