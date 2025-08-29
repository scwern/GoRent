package service

import (
	"GoRent/internal/repository/analytics"
	"context"
	"fmt"
	"time"
)

type AnalyticsService interface {
	GetProfit(ctx context.Context, fromDate, toDate time.Time) (float64, error)
	GetPopularBrands(ctx context.Context, fromDate, toDate time.Time, limit int) ([]analytics.BrandRevenue, error)
	GetRentalStats(ctx context.Context, fromDate, toDate time.Time) (*analytics.RentalStats, error)
}

type analyticsService struct {
	analyticsRepo analytics.Repository
}

func NewAnalyticsService(analyticsRepo analytics.Repository) AnalyticsService {
	return &analyticsService{analyticsRepo: analyticsRepo}
}

func (s *analyticsService) GetProfit(ctx context.Context, fromDate, toDate time.Time) (float64, error) {
	if fromDate.After(toDate) {
		return 0, fmt.Errorf("fromDate cannot be after toDate")
	}

	return s.analyticsRepo.GetProfitAnalytics(ctx, fromDate, toDate)
}

func (s *analyticsService) GetPopularBrands(ctx context.Context, fromDate, toDate time.Time, limit int) ([]analytics.BrandRevenue, error) {
	if fromDate.After(toDate) {
		return nil, fmt.Errorf("fromDate cannot be after toDate")
	}

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	return s.analyticsRepo.GetPopularBrands(ctx, fromDate, toDate, limit)
}

func (s *analyticsService) GetRentalStats(ctx context.Context, fromDate, toDate time.Time) (*analytics.RentalStats, error) {
	if fromDate.After(toDate) {
		return nil, fmt.Errorf("fromDate cannot be after toDate")
	}

	return s.analyticsRepo.GetRentalStats(ctx, fromDate, toDate)
}
