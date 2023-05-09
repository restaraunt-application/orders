package db

import (
	"context"
	"orders/internal/infra/config"

	"github.com/uptrace/opentelemetry-go-extra/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"

	_ "github.com/jackc/pgx/v4/stdlib" // import for pgx support
	"github.com/jmoiron/sqlx"
	"github.com/uptrace/opentelemetry-go-extra/otelsqlx"
)

func NewDB(ctx context.Context, dbConfig *config.PostgreSQLConfig) (*sqlx.DB, error) {
	db, err := otelsqlx.ConnectContext(ctx, "pgx", dbConfig.ConnectionString(),
		otelsql.WithAttributes(semconv.DBSystemPostgreSQL), otelsql.WithDBName(dbConfig.DBName))
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(dbConfig.MaxOpenConnections)

	return db, nil
}
