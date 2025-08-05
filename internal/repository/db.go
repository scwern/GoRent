package repository

import (
	"GoRent/internal/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgresDB(cfg *config.Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres",
		"host="+cfg.DB.Host+
			" port="+cfg.DB.Port+
			" user="+cfg.DB.User+
			" password="+cfg.DB.Password+
			" dbname="+cfg.DB.Name+
			" sslmode="+cfg.DB.SSLMode)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
