package main

import (
	"github.com/carinfinin/loyalty-program/config"
	"github.com/carinfinin/loyalty-program/internal/server"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	s := server.New(&config.Config{
		Addr:         ":8080",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	})

	go func() {
		s.Run()
	}()

	<-exit
	// todo stoping

}
