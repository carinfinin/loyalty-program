package main

import (
	"context"
	"github.com/carinfinin/loyalty-program/internal/config"
	"github.com/carinfinin/loyalty-program/internal/logger"
	"github.com/carinfinin/loyalty-program/internal/server"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
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
	migration(cfg)

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

	ctx, cancel := context.WithTimeout(context.Background(), cfg.WriteTimeout)
	defer cancel()

	if err = s.Shutdown(ctx); err != nil {
		logger.Log.Info("stop shutdown error")
		s.Close()
	}

	logger.Log.Info("stop app")
}

func migration(cfg *config.Config) {
	m, err := migrate.New(
		"file://migrations",
		cfg.DBPath)
	if err != nil {
		panic(err)
	}
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return
		}
		panic(err)
	}
}
