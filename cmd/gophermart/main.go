package main

import (
	"github.com/carinfinin/loyalty-program/internal/config"
	"github.com/carinfinin/loyalty-program/internal/logger"
	"github.com/carinfinin/loyalty-program/internal/server"
	"os"
	"os/signal"
	"syscall"
)

func main() {

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	cfg, err := config.New()
	if err != nil {
		panic(err)
	}
	err = logger.Configure(cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	s := server.New(cfg)

	go func() {
		s.Run()
		logger.Log.Info("start app")

	}()

	<-exit

	logger.Log.Info("stop app")
	// todo stoping
}
