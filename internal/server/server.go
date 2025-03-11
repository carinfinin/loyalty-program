package server

import (
	"github.com/carinfinin/loyalty-program/internal/config"
	"github.com/carinfinin/loyalty-program/internal/logger"
	"github.com/carinfinin/loyalty-program/internal/router"
	"github.com/carinfinin/loyalty-program/internal/service"
	"github.com/carinfinin/loyalty-program/internal/store/storepg"
	"net/http"
)

type Server struct {
	*http.Server
}

func New(cfg *config.Config) *Server {
	store, err := storepg.New(cfg)
	if err != nil {
		logger.Log.Error("run store error: ", err)
		panic(err)
	}
	service := service.New(cfg, store)
	r := router.Configure(cfg, service)
	return &Server{
		Server: &http.Server{
			Handler:      r.Handler,
			Addr:         cfg.Addr,
			WriteTimeout: cfg.WriteTimeout,
			ReadTimeout:  cfg.ReadTimeout,
		},
	}
}

func (s *Server) Run() error {
	logger.Log.Info("start server on ", s.Addr)
	return s.ListenAndServe()
}
