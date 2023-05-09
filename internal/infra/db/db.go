package db

import (
	"orders/internal/infra/config"

	_ "github.com/jackc/pgx/v4/stdlib" // import for pgx support
	"github.com/jmoiron/sqlx"
)

func NewDB(config *config.PostgreSQLConfig) (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", config.ConnectionString())
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(config.MaxOpenConnections)

	return db, nil
}
