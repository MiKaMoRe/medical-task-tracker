package config

import (
	"fmt"
	"os"
	"strings"
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

func loadEnv(envErrs *[]string, key string) string {
	env := os.Getenv(key)

	if len(env) == 0 {
		*envErrs = append(*envErrs, key)
	}

	return env
}

func validateAllowed(value, name string, allowed ...string) string {
	for _, a := range allowed {
		if value == a {
			return ""
		}
	}
	return fmt.Sprintf("%s must be one of %s, got %q", name, strings.Join(allowed, ", "), value)
}

func NewConfig() (*Config, error) {
	envErrs := []string{}
	env := loadEnv(&envErrs, "APP_ENV")

	logLevel := loadEnv(&envErrs, "LOG_LEVEL")

	dbHost := loadEnv(&envErrs, "POSTGRES_HOST")

	dbConfig := &DBConfig{
		Host: dbHost,
		Port: loadEnv(&envErrs, "POSTGRES_PORT"),
		User: loadEnv(&envErrs, "POSTGRES_USER"),
		Pass: loadEnv(&envErrs, "POSTGRES_PASSWORD"),
		Db:   loadEnv(&envErrs, "POSTGRES_DB"),
	}

	if len(envErrs) > 0 {
		return nil, fmt.Errorf("missing env variables: %s", strings.Join(envErrs, ", "))
	}

	validationErrs := []string{}
	if msg := validateAllowed(env, "APP_ENV", EnvDev, EnvTest, EnvProd); msg != "" {
		validationErrs = append(validationErrs, msg)
	}
	if msg := validateAllowed(logLevel, "LOG_LEVEL", LogLevelDebug, LogLevelInfo, LogLevelWarn, LogLevelError); msg != "" {
		validationErrs = append(validationErrs, msg)
	}
	if len(validationErrs) > 0 {
		return nil, fmt.Errorf("invalid config: %s", strings.Join(validationErrs, "; "))
	}

	return &Config{
		Env:      env,
		LogLevel: logLevel,
		AppPort:  loadEnv(&envErrs, "APP_PORT"),

		DB: dbConfig,
	}, nil
}
