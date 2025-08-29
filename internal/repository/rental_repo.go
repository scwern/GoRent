package repository

import (
	"GoRent/internal/domain/rental"
	"GoRent/internal/repository/db"
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type RentalRepository interface {
	Create(ctx context.Context, rental *rental.Rental) error
	GetByID(ctx context.Context, id uuid.UUID) (*rental.Rental, error)
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]*rental.Rental, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status rental.Status) error
	CheckCarAvailability(ctx context.Context, carID uuid.UUID, startDate, endDate time.Time) (bool, error)
}

type rentalRepo struct {
	db *db.DB
}

func NewRentalRepository(db *db.DB) RentalRepository {
	return &rentalRepo{db: db}
}

func (r *rentalRepo) Create(ctx context.Context, rental *rental.Rental) error {
	query := `
        INSERT INTO rentals (id, car_id, user_id, start_date, end_date, total_price, status)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

	_, err := r.db.ExecContext(ctx, query,
		rental.ID,
		rental.CarID,
		rental.UserID,
		rental.StartDate,
		rental.EndDate,
		rental.TotalPrice,
		rental.Status,
	)

	if err != nil {
		return fmt.Errorf("failed to create rental: %w", err)
	}

	return nil
}

func (r *rentalRepo) GetByID(ctx context.Context, id uuid.UUID) (*rental.Rental, error) {
	query := `
		SELECT id, car_id, user_id, start_date, end_date, total_price, status, created_at
		FROM rentals WHERE id = $1
	`

	var rent rental.Rental
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&rent.ID,
		&rent.CarID,
		&rent.UserID,
		&rent.StartDate,
		&rent.EndDate,
		&rent.TotalPrice,
		&rent.Status,
		&rent.CreatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to get rental: %w", err)
	}

	return &rent, nil
}

func (r *rentalRepo) GetByUserID(ctx context.Context, userID uuid.UUID) ([]*rental.Rental, error) {
	query := `
		SELECT id, car_id, user_id, start_date, end_date, total_price, status, created_at
		FROM rentals WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to query rentals: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			fmt.Printf("warning: failed to close rows: %v\n", err)
		}
	}()

	var rentals []*rental.Rental
	for rows.Next() {
		var rent rental.Rental
		if err := rows.Scan(
			&rent.ID,
			&rent.CarID,
			&rent.UserID,
			&rent.StartDate,
			&rent.EndDate,
			&rent.TotalPrice,
			&rent.Status,
			&rent.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("failed to scan rental: %w", err)
		}
		rentals = append(rentals, &rent)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return rentals, nil
}

func (r *rentalRepo) UpdateStatus(ctx context.Context, id uuid.UUID, status rental.Status) error {
	query := `UPDATE rentals SET status = $1 WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update rental status: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("rental not found")
	}

	return nil
}

func (r *rentalRepo) CheckCarAvailability(ctx context.Context, carID uuid.UUID, startDate, endDate time.Time) (bool, error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return false, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	lockQuery := `SELECT id FROM cars WHERE id = $1 FOR UPDATE`
	_, err = tx.ExecContext(ctx, lockQuery, carID)
	if err != nil {
		return false, fmt.Errorf("failed to lock car: %w", err)
	}

	availabilityQuery := `
        SELECT COUNT(*) FROM rentals 
        WHERE car_id = $1 
        AND status IN ('pending', 'active')
        AND (start_date, end_date) OVERLAPS ($2, $3)
    `

	var count int
	err = tx.QueryRowContext(ctx, availabilityQuery, carID, startDate, endDate).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check availability: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return false, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return count == 0, nil
}
