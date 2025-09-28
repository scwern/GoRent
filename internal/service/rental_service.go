package service

import (
	_ "GoRent/internal/domain/car"
	"GoRent/internal/domain/rental"
	"GoRent/internal/repository"
	"context"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type RentalService interface {
	CreateRental(ctx context.Context, input rental.CreateRequest, userID uuid.UUID) (*rental.Response, error)
	GetUserRentals(ctx context.Context, userID uuid.UUID) ([]rental.Response, error)
	GetRental(ctx context.Context, rentalID uuid.UUID, userID uuid.UUID) (*rental.Response, error)
	CancelRental(ctx context.Context, rentalID uuid.UUID, userID uuid.UUID) error
	ApproveRental(ctx context.Context, rentalID uuid.UUID) error
}

type rentalService struct {
	rentalRepo repository.RentalRepository
	carRepo    repository.CarRepository
}

func NewRentalService(rentalRepo repository.RentalRepository, carRepo repository.CarRepository) RentalService {
	return &rentalService{
		rentalRepo: rentalRepo,
		carRepo:    carRepo,
	}
}

func (s *rentalService) CreateRental(ctx context.Context, input rental.CreateRequest, userID uuid.UUID) (*rental.Response, error) {
	if err := validateRentalDates(input.StartDate, input.EndDate); err != nil {
		return nil, err
	}

	car, err := s.carRepo.GetByID(ctx, input.CarID)
	if err != nil {
		return nil, fmt.Errorf("car not found: %w", err)
	}

	if !car.IsAvailable {
		return nil, fmt.Errorf("car is not available for rental")
	}

	available, err := s.rentalRepo.CheckCarAvailability(ctx, input.CarID, input.StartDate, input.EndDate)
	if err != nil {
		return nil, fmt.Errorf("availability check failed: %w", err)
	}
	if !available {
		return nil, fmt.Errorf("car is not available for selected dates")
	}

	days := calculateRentalDays(input.StartDate, input.EndDate)
	totalPrice := car.PricePerDay * float64(days)

	newRental := &rental.Rental{
		ID:         uuid.New(),
		CarID:      input.CarID,
		UserID:     userID,
		StartDate:  input.StartDate,
		EndDate:    input.EndDate,
		TotalPrice: totalPrice,
		Status:     rental.StatusPending,
		CreatedAt:  time.Now(),
	}

	if err := s.rentalRepo.Create(ctx, newRental); err != nil {
		return nil, fmt.Errorf("failed to create rental: %w", err)
	}

	return s.entityToResponse(newRental), nil
}

func (s *rentalService) GetUserRentals(ctx context.Context, userID uuid.UUID) ([]rental.Response, error) {
	rentals, err := s.rentalRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user rentals: %w", err)
	}

	responses := make([]rental.Response, 0, len(rentals))
	for _, r := range rentals {
		responses = append(responses, *s.entityToResponse(r))
	}

	return responses, nil
}

func (s *rentalService) GetRental(ctx context.Context, rentalID uuid.UUID, userID uuid.UUID) (*rental.Response, error) {
	rentalItem, err := s.rentalRepo.GetByID(ctx, rentalID)
	if err != nil {
		return nil, fmt.Errorf("failed to get rental: %w", err)
	}

	if rentalItem.UserID != userID {
		return nil, fmt.Errorf("access denied")
	}

	return s.entityToResponse(rentalItem), nil
}

func (s *rentalService) CancelRental(ctx context.Context, rentalID uuid.UUID, userID uuid.UUID) error {
	rentalItem, err := s.rentalRepo.GetByID(ctx, rentalID)
	if err != nil {
		return fmt.Errorf("failed to get rental: %w", err)
	}

	if rentalItem.UserID != userID {
		return fmt.Errorf("access denied")
	}

	if rentalItem.Status != rental.StatusPending {
		return fmt.Errorf("only pending rentals can be canceled")
	}

	if err := s.rentalRepo.UpdateStatus(ctx, rentalID, rental.StatusCanceled); err != nil {
		return fmt.Errorf("failed to cancel rental: %w", err)
	}

	return nil
}

func validateRentalDates(startDate, endDate time.Time) error {
	if startDate.Before(time.Now().AddDate(0, 0, -1)) {
		return fmt.Errorf("start date cannot be in the past")
	}
	if endDate.Before(startDate) {
		return fmt.Errorf("end date must be after start date")
	}
	if endDate.Sub(startDate).Hours() < 24 {
		return fmt.Errorf("minimum rental period is 24 hours")
	}
	return nil
}

func calculateRentalDays(startDate, endDate time.Time) int {
	days := int(endDate.Sub(startDate).Hours() / 24)
	if days < 1 {
		return 1
	}
	return days
}

func (s *rentalService) entityToResponse(r *rental.Rental) *rental.Response {
	return &rental.Response{
		ID:         r.ID.String(),
		CarID:      r.CarID.String(),
		UserID:     r.UserID.String(),
		StartDate:  r.StartDate,
		EndDate:    r.EndDate,
		TotalPrice: r.TotalPrice,
		Status:     r.Status,
		CreatedAt:  r.CreatedAt,
	}
}

func (s *rentalService) ApproveRental(ctx context.Context, rentalID uuid.UUID) error {
	rentalItem, err := s.rentalRepo.GetByID(ctx, rentalID)
	if err != nil {
		return fmt.Errorf("failed to get rental: %w", err)
	}

	if rentalItem.Status != rental.StatusPending {
		return fmt.Errorf("only pending rentals can be approved")
	}

	if err := s.rentalRepo.UpdateStatus(ctx, rentalID, rental.StatusActive); err != nil {
		return fmt.Errorf("failed to approve rental: %w", err)
	}

	return nil
}
