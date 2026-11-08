package configuration

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type DatabaseConfiguration struct {
	Host        string
	Database    string
	User        string
	Password    string
	ExtraParams string
}

var (
	ErrDbHostnameNotSet = errors.New("Failed to read DB__HOSTNAME from environment")
	ErrDbDatabaseNotSet = errors.New("Failed to read DB__DATABASE from environment")
	ErrDbUserNotSet     = errors.New("Failed to read DB__USER from environment")
	ErrDbPasswordNotSet = errors.New("Failed to read DB__PASSWORD from environment")
)

func ReadDatabaseConfigurationFromEnv() (DatabaseConfiguration, error) {
	host := os.Getenv("DB__HOSTNAME")
	if host == "" {
		return DatabaseConfiguration{}, ErrDbHostnameNotSet
	}

	database := os.Getenv("DB__DATABASE")
	if database == "" {
		return DatabaseConfiguration{}, ErrDbDatabaseNotSet
	}

	user := os.Getenv("DB__USER")
	if user == "" {
		return DatabaseConfiguration{}, ErrDbUserNotSet
	}

	password := os.Getenv("DB__PASSWORD")
	if password == "" {
		return DatabaseConfiguration{}, ErrDbPasswordNotSet
	}

	extraParams := os.Getenv("DB__EXTRAPARAMS")

	return DatabaseConfiguration{
		Host:        host,
		Database:    database,
		User:        user,
		Password:    password,
		ExtraParams: extraParams,
	}, nil
}

func (c *DatabaseConfiguration) CreateDatabase(ctx context.Context) (*pgxpool.Pool, error) {
	connectionString := url.URL{
		Scheme:   "postgres",
		Host:     c.Host,
		Path:     c.Database,
		User:     url.UserPassword(c.User, c.Password),
		RawQuery: c.ExtraParams,
	}

	pool, err := pgxpool.New(ctx, connectionString.String())
	if err != nil {
		return nil, fmt.Errorf("Failed to initialise database connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("Failed to ping new database connection pool: %w", err)
	}

	return pool, nil
}
