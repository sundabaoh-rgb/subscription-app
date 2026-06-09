package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App      AppConfig
	Database DatabaseConfig
}

type AppConfig struct {
	Port     string `env:"APP_PORT"`
	LogLevel string `env:"LOG_LEVEL"     env-default:"info"`
	Env      string `env:"APP_ENV"       env-default:"local"`
}

type DatabaseConfig struct {
	Host     string `env:"DB_HOST"     env-required:"true"`
	Port     string `env:"DB_PORT"     env-default:"5432"`
	Name     string `env:"DB_NAME"     env-required:"true"`
	User     string `env:"DB_USER"     env-required:"true"`
	Password string `env:"DB_PASSWORD" env-required:"true"`
}

func MustLoad() *Config {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("config error: %s", err)
	}
	return &cfg
}
