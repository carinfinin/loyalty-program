package config

import (
	"flag"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type Config struct {
	Addr         string        `yaml:"addr" env:"RUN_ADDRESS" env-default:"localhost:8080"`
	WriteTimeout time.Duration `yaml:"timeout_write" env:"TIMEOUT_WRITE" env-default:"15s"`
	ReadTimeout  time.Duration `yaml:"timeout_read" env:"TIMEOUT_READ" env-default:"15s"`
	DBPath       string        `yaml:"db_path" env:"DATABASE_URI" env-default:"postgres://user:password@localhost:5432/loyalty?sslmode=disable"`
	AccrualAddr  string        `yaml:"accrual_address" env:"ACCRUAL_SYSTEM_ADDRESS" env-default:""`
	LogLevel     string        `yaml:"log_level" env:"LOG_LEVEL" env-default:"info"`
}

func New() (*Config, error) {
	var cfg = Config{}

	err := cleanenv.ReadConfig("config/config.yml", &cfg)
	if err != nil {
		return nil, err
	}

	var runAddr string
	var dbPath string
	var accrualAddr string

	flag.StringVar(&runAddr, "a", "", "RUN_ADDRESS")
	flag.StringVar(&dbPath, "b", "", "DATABASE_URI")
	flag.StringVar(&accrualAddr, "r", "", "ACCRUAL_SYSTEM_ADDRESS")

	if runAddr != "" {
		cfg.Addr = runAddr
	}
	if runAddr != "" {
		cfg.DBPath = dbPath
	}
	if runAddr != "" {
		cfg.AccrualAddr = accrualAddr
	}
	return &cfg, nil
}
