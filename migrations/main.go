package main

import (
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"net/url"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed sql/*.sql
var embeddedMigrations embed.FS

func main() {
	connectionString, err := readConnectionStringFromEnv()
	if err != nil {
		panic(fmt.Errorf("Failed to load database configuration: %w", err))
	}

	db, err := sql.Open("pgx", connectionString)
	if err != nil {
		panic(fmt.Errorf("Failed to create database: %w", err))
	}
	defer db.Close()

	goose.SetBaseFS(embeddedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		panic(fmt.Errorf("Failed to set postgres dialect for goose: %w", err))
	}

	if err := goose.Up(db, "migrations"); err != nil {
		panic(fmt.Errorf("Failed to run migrations: %w", err))
	}
}

func readConnectionStringFromEnv() (string, error) {
	host := os.Getenv("DB__HOSTNAME")
	if host == "" {
		return "", errors.New("Failed to read DB__HOSTNAME from environment variables")
	}

	database := os.Getenv("DB__DATABASE")
	if database == "" {
		return "", errors.New("Failed to read DB__DATABASE from environment variables")
	}

	user := os.Getenv("DB__USER")
	if user == "" {
		return "", errors.New("Failed to read DB__USER from environment variables")
	}

	password := os.Getenv("DB__PASSWORD")
	if password == "" {
		return "", errors.New("Failed to read DB__PASSWORD from environment variables")
	}

	extraParams := os.Getenv("DB__EXTRAPARAMS")

	connectionStringUrl := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(user, password),
		Host:     host,
		Path:     database,
		RawQuery: extraParams,
	}

	return connectionStringUrl.String(), nil
}
