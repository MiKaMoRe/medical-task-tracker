package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

const (
	EnvDev  = "dev"
	EnvTest = "test"
	EnvProd = "prod"
)

const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
)

type Config struct {
	Env      string
	LogLevel string
	AppPort  string

	DB *DBConfig
}

type DBConfig struct {
	Host string
	Port string
	User string
	Pass string
	Db   string
}

func (c *DBConfig) PostgresDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		c.Host,
		c.User,
		c.Pass,
		c.Db,
		c.Port,
	)
}

func mustLoad(key string) string {
	env := os.Getenv(key)

	if len(env) == 0 {
		panic("Missing Env variable: " + key)
	}

	return env
}

func NewConfig() (*Config, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = EnvProd
	}

	if env == EnvDev {
		if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("loading .env: %w", err)
		}
	}

	logLevel := os.Getenv("LOG_LEVEL")
	if len(logLevel) == 0 {
		logLevel = LogLevelInfo
	}

	dbHost := os.Getenv("POSTGRES_HOST")

	dbConfig := &DBConfig{
		Host: dbHost,
		Port: mustLoad("POSTGRES_PORT"),
		User: mustLoad("POSTGRES_USER"),
		Pass: mustLoad("POSTGRES_PASSWORD"),
		Db:   mustLoad("POSTGRES_DB"),
	}

	return &Config{
		Env:      env,
		LogLevel: logLevel,
		AppPort:  mustLoad("APP_PORT"),

		DB: dbConfig,
	}, nil
}
