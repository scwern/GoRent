package analytics

import (
	"GoRent/internal/repository/db"
	"context"
	"fmt"
	"time"
)

type Repository interface {
	GetProfitAnalytics(ctx context.Context, fromDate, toDate time.Time) (float64, error)
	GetPopularBrands(ctx context.Context, fromDate, toDate time.Time, limit int) ([]BrandRevenue, error)
	GetRentalStats(ctx context.Context, fromDate, toDate time.Time) (*RentalStats, error)
}

type analyticsRepo struct {
	db *db.DB
}

func NewRepository(db *db.DB) Repository {
	return &analyticsRepo{db: db}
}

func (r *analyticsRepo) GetProfitAnalytics(ctx context.Context, fromDate, toDate time.Time) (float64, error) {
	query := `
		SELECT COALESCE(SUM(total_price), 0) 
		FROM rentals 
		WHERE status IN ('active', 'completed')
		AND created_at BETWEEN $1 AND $2
	`

	var totalProfit float64
	err := r.db.QueryRowContext(ctx, query, fromDate, toDate).Scan(&totalProfit)
	if err != nil {
		return 0, fmt.Errorf("failed to get profit analytics: %w", err)
	}

	return totalProfit, nil
}

func (r *analyticsRepo) GetPopularBrands(ctx context.Context, fromDate, toDate time.Time, limit int) ([]BrandRevenue, error) {
	query := `
		SELECT c.brand, 
			   SUM(r.total_price) as revenue,
			   COUNT(r.id) as rental_count
		FROM rentals r
		JOIN cars c ON r.car_id = c.id
		WHERE r.status IN ('active', 'completed')
		AND r.created_at BETWEEN $1 AND $2
		GROUP BY c.brand
		ORDER BY revenue DESC
		LIMIT $3
	`

	rows, err := r.db.QueryContext(ctx, query, fromDate, toDate, limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get popular brands: %w", err)
	}
	defer rows.Close()

	var brands []BrandRevenue
	for rows.Next() {
		var br BrandRevenue
		if err := rows.Scan(&br.Brand, &br.Revenue, &br.Count); err != nil {
			return nil, fmt.Errorf("failed to scan brand revenue: %w", err)
		}
		brands = append(brands, br)
	}

	return brands, nil
}

func (r *analyticsRepo) GetRentalStats(ctx context.Context, fromDate, toDate time.Time) (*RentalStats, error) {
	query := `
		SELECT 
			COUNT(*) as total_rentals,
			COALESCE(SUM(total_price), 0) as total_revenue,
			COALESCE(AVG(EXTRACT(EPOCH FROM (end_date - start_date)) / 86400), 0) as avg_rental_days,
			COALESCE(SUM(total_price) / NULLIF(COUNT(*), 0), 0) as avg_daily_revenue
		FROM rentals 
		WHERE status IN ('active', 'completed')
		AND created_at BETWEEN $1 AND $2
	`

	var stats RentalStats
	err := r.db.QueryRowContext(ctx, query, fromDate, toDate).Scan(
		&stats.TotalRentals,
		&stats.TotalRevenue,
		&stats.AvgRentalDays,
		&stats.AvgDailyRevenue,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get rental stats: %w", err)
	}

	return &stats, nil
}
