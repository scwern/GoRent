package repository

import (
	"GoRent/internal/domain/car"
	"GoRent/internal/repository/db"
	"context"
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"
)

type CarRepository interface {
	Create(ctx context.Context, c *car.Car) error
	GetByID(ctx context.Context, id uuid.UUID) (*car.Car, error)
	GetAll(ctx context.Context, filters map[string]interface{}) ([]*car.Car, error)
	Update(ctx context.Context, id uuid.UUID, updateData *car.UpdateRequest) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type carRepo struct {
	db *db.DB
}

func NewCarRepository(db *db.DB) CarRepository {
	return &carRepo{db: db}
}

func (r *carRepo) Create(ctx context.Context, c *car.Car) error {
	const query = `
		INSERT INTO cars (id, model, brand, year, price_per_day, is_available)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		c.ID,
		c.Model,
		c.Brand,
		c.Year,
		c.PricePerDay,
		c.IsAvailable,
	)
	if err != nil {
		return fmt.Errorf("failed to create car: %w", err)
	}
	return nil
}

func (r *carRepo) GetByID(ctx context.Context, id uuid.UUID) (*car.Car, error) {
	const query = `
		SELECT id, model, brand, year, price_per_day, is_available 
		FROM cars 
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)

	var c car.Car
	err := row.Scan(
		&c.ID,
		&c.Model,
		&c.Brand,
		&c.Year,
		&c.PricePerDay,
		&c.IsAvailable,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get car by ID %s: %w", id, err)
	}

	return &c, nil
}

func (r *carRepo) GetAll(ctx context.Context, filters map[string]interface{}) ([]*car.Car, error) {
	var (
		query strings.Builder
		args  []interface{}
	)

	query.WriteString(`
		SELECT id, model, brand, year, price_per_day, is_available 
		FROM cars 
		WHERE 1=1
	`)

	argPos := 1

	if brand, ok := filters["brand"].(string); ok && brand != "" {
		query.WriteString(fmt.Sprintf(" AND brand = $%d", argPos))
		args = append(args, brand)
		argPos++
	}

	if minPrice, ok := filters["min_price"].(float64); ok {
		query.WriteString(fmt.Sprintf(" AND price_per_day >= $%d", argPos))
		args = append(args, minPrice)
		argPos++
	}

	if maxPrice, ok := filters["max_price"].(float64); ok {
		query.WriteString(fmt.Sprintf(" AND price_per_day <= $%d", argPos))
		args = append(args, maxPrice)
		argPos++
	}

	if available, ok := filters["available"].(bool); ok {
		query.WriteString(fmt.Sprintf(" AND is_available = $%d", argPos))
		args = append(args, available)
		argPos++
	}

	rows, err := r.db.QueryContext(ctx, query.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query cars: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}()

	var cars []*car.Car
	for rows.Next() {
		var c car.Car
		if err := rows.Scan(
			&c.ID,
			&c.Model,
			&c.Brand,
			&c.Year,
			&c.PricePerDay,
			&c.IsAvailable,
		); err != nil {
			return nil, fmt.Errorf("failed to scan car: %w", err)
		}
		cars = append(cars, &c)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return cars, nil
}

func (r *carRepo) Update(ctx context.Context, id uuid.UUID, updateData *car.UpdateRequest) error {
	var (
		query strings.Builder
		args  []interface{}
	)

	query.WriteString("UPDATE cars SET ")
	argPos := 1

	if updateData.Model != nil {
		query.WriteString(fmt.Sprintf("model = $%d, ", argPos))
		args = append(args, *updateData.Model)
		argPos++
	}

	if updateData.Brand != nil {
		query.WriteString(fmt.Sprintf("brand = $%d, ", argPos))
		args = append(args, *updateData.Brand)
		argPos++
	}

	if updateData.Year != nil {
		query.WriteString(fmt.Sprintf("year = $%d, ", argPos))
		args = append(args, *updateData.Year)
		argPos++
	}

	if updateData.PricePerDay != nil {
		query.WriteString(fmt.Sprintf("price_per_day = $%d, ", argPos))
		args = append(args, *updateData.PricePerDay)
		argPos++
	}

	if updateData.IsAvailable != nil {
		query.WriteString(fmt.Sprintf("is_available = $%d, ", argPos))
		args = append(args, *updateData.IsAvailable)
		argPos++
	}

	if argPos == 1 {
		return errors.New("no fields to update")
	}

	queryStr := strings.TrimSuffix(query.String(), ", ")
	queryStr += fmt.Sprintf(" WHERE id = $%d", argPos)
	args = append(args, id)

	result, err := r.db.ExecContext(ctx, queryStr, args...)
	if err != nil {
		return fmt.Errorf("failed to update car: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("car not found")
	}

	return nil
}

func (r *carRepo) Delete(ctx context.Context, id uuid.UUID) error {
	const query = "DELETE FROM cars WHERE id = $1"
	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("failed to delete car: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return errors.New("car not found")
	}

	return nil
}
