package server

import (
	"github.com/carinfinin/loyalty-program/config"
	"github.com/carinfinin/loyalty-program/internal/logger"
	"github.com/carinfinin/loyalty-program/internal/router"
	"net/http"
)

type Server struct {
	*http.Server
}

func New(cfg *config.Config) *Server {
	r := router.Configure(cfg)
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
