package config

import "time"

type Config struct {
	Addr         string
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
}
