package utils

import (
	"GoRent/internal/repository/db"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

func SetupTestDB() (*db.DB, func(), error) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "postgres:15",
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "test",
		},
		WaitingFor: wait.ForLog("database system is ready to accept connections"),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start container: %w", err)
	}

	host, err := container.Host(ctx)
	if err != nil {
		container.Terminate(ctx)
		return nil, nil, fmt.Errorf("failed to get host: %w", err)
	}

	port, err := container.MappedPort(ctx, "5432")
	if err != nil {
		container.Terminate(ctx)
		return nil, nil, fmt.Errorf("failed to get port: %w", err)
	}

	connStr := fmt.Sprintf("host=%s port=%s user=test password=test dbname=test sslmode=disable",
		host, port.Port())

	sqlxDB, err := sqlx.Open("postgres", connStr)
	if err != nil {
		container.Terminate(ctx)
		return nil, nil, fmt.Errorf("failed to open DB: %w", err)
	}

	if err := waitForDB(sqlxDB); err != nil {
		sqlxDB.Close()
		container.Terminate(ctx)
		return nil, nil, fmt.Errorf("failed to wait for DB: %w", err)
	}

	if err := applySchema(sqlxDB); err != nil {
		sqlxDB.Close()
		container.Terminate(ctx)
		return nil, nil, fmt.Errorf("failed to apply schema: %w", err)
	}

	cleanup := func() {
		if err := sqlxDB.Close(); err != nil {
			log.Printf("Failed to close DB: %v", err)
		}
		if err := container.Terminate(ctx); err != nil {
			log.Printf("Failed to terminate container: %v", err)
		}
	}

	return db.NewDB(sqlxDB), cleanup, nil
}

func waitForDB(db *sqlx.DB) error {
	for i := 0; i < 10; i++ {
		err := db.Ping()
		if err == nil {
			return nil
		}
		time.Sleep(1 * time.Second)
	}
	return fmt.Errorf("failed to connect to DB")
}

func applySchema(db *sqlx.DB) error {
	schema := `
		DROP TABLE IF EXISTS rentals CASCADE;
		DROP TABLE IF EXISTS cars CASCADE;
		DROP TABLE IF EXISTS users CASCADE;
		DROP TYPE IF EXISTS user_role CASCADE;
		DROP TYPE IF EXISTS rental_status CASCADE;

		CREATE TYPE user_role AS ENUM ('admin', 'manager', 'client');
		CREATE TYPE rental_status AS ENUM ('pending', 'active', 'completed', 'canceled');

		CREATE TABLE users (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			name VARCHAR(255) NOT NULL,
			email VARCHAR(255) UNIQUE NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			role user_role NOT NULL DEFAULT 'client'
		);

		CREATE TABLE cars (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			model VARCHAR(255) NOT NULL,
			brand VARCHAR(255) NOT NULL,
			year INTEGER NOT NULL,
			price_per_day NUMERIC(10,2) NOT NULL,
			is_available BOOLEAN NOT NULL DEFAULT TRUE
		);

		CREATE TABLE rentals (
			id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
			car_id UUID NOT NULL REFERENCES cars(id),
			user_id UUID NOT NULL REFERENCES users(id),
			start_date TIMESTAMP NOT NULL,
			end_date TIMESTAMP NOT NULL,
			total_price NUMERIC(10,2) NOT NULL,
			status rental_status NOT NULL DEFAULT 'pending',
			created_at TIMESTAMP NOT NULL DEFAULT NOW()
		);
	`

	_, err := db.Exec(schema)
	return err
}

func CreateTestCar(db *db.DB, id, model, brand string, year int, pricePerDay float64, available bool) error {
	query := `INSERT INTO cars (id, model, brand, year, price_per_day, is_available) VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := db.ExecContext(context.Background(), query, id, model, brand, year, pricePerDay, available)
	return err
}

func CreateTestUser(db *db.DB, id, email, role string) error {
	query := `INSERT INTO users (id, name, email, password_hash, role) VALUES ($1, $2, $3, $4, $5)`
	passwordHash := "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi"
	_, err := db.ExecContext(context.Background(), query, id, "Test User", email, passwordHash, role)
	return err
}
