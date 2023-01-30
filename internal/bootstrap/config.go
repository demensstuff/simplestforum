package bootstrap

import (
	"errors"
	"os"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

const defaultEnvPath = "deployment/.env"

// DBConfig contains all the database configuration info.
type DBConfig struct {
	Scheme   string `envconfig:"DB_SCHEME"`
	Host     string `envconfig:"DB_HOST"`
	Port     string `envconfig:"DB_PORT"`
	Name     string `envconfig:"DB_NAME"`
	Username string `envconfig:"DB_USERNAME"`
	Password string `envconfig:"DB_PASSWORD"`
}

// Config contains all the configuration info.
type Config struct {
	HTTPPort string `envconfig:"HTTP_PORT"`
	DB       DBConfig
}

// NewConfig loads configuration from the environment variables, optionally loading them from the file.
func NewConfig() (*Config, error) {
	err := godotenv.Load(defaultEnvPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}

	var cfg Config

	err = envconfig.Process("", &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
