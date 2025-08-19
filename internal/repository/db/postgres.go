package db

import (
	"GoRent/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgresDB(cfg *config.Config) (*DB, error) {
	connStr := "host=" + cfg.DB.Host +
		" port=" + cfg.DB.Port +
		" user=" + cfg.DB.User +
		" password=" + cfg.DB.Password +
		" dbname=" + cfg.DB.Name +
		" sslmode=" + cfg.DB.SSLMode

	sqlxDB, err := sqlx.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := sqlxDB.Ping(); err != nil {
		return nil, err
	}

	return NewDB(sqlxDB), nil
}
