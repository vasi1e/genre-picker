package env

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type DBConfig struct {
	Name     string `envconfig:"DB_NAME" required:"true"`
	Username string `envconfig:"DB_USERNAME" required:"true"`
	Password string `envconfig:"DB_PASSWORD" required:"true"`
	Host     string `envconfig:"DB_HOST" required:"true"`
}

func LoadDBConfig() (DBConfig, error) {
	var dbConfig DBConfig
	if err := envconfig.Process("", &dbConfig); err != nil {
		return DBConfig{}, fmt.Errorf("failed load DB config from environment: %w", err)
	}
	return dbConfig, nil
}
