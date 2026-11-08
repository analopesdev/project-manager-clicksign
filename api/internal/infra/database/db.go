package database

import (
	"database/sql"
	"embed"
	"fmt"

	"github.com/analopesdev/project-manager-clicksign/api/internal/infra/config"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrations embed.FS

func NewPostgresDB(cfg *config.Config) (*sql.DB, error) {
	connConfig, err := pgx.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, fmt.Errorf("parsing database url: %w", err)
	}
	connConfig.RuntimeParams["timezone"] = "America/Sao_Paulo"

	db := stdlib.OpenDB(*connConfig)

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("connecting to database: %w", err)
	}

	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("postgres"); err != nil {
		return nil, err
	}
	if err := goose.Up(db, "migrations"); err != nil {
		return nil, fmt.Errorf("running migrations: %w", err)
	}

	return db, nil
}
