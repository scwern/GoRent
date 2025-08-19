package service

import (
	"GoRent/internal/domain/car"
	"GoRent/internal/repository"
	"context"
	"fmt"
	"github.com/google/uuid"
)

type CarService interface {
	CreateCar(ctx context.Context, input car.CreateRequest) (*car.Response, error)
	GetCar(ctx context.Context, id uuid.UUID) (*car.Response, error)
	GetAllCars(ctx context.Context, filters map[string]interface{}) ([]car.Response, error)
	UpdateCar(ctx context.Context, id uuid.UUID, input car.UpdateRequest) (*car.Response, error)
	DeleteCar(ctx context.Context, id uuid.UUID) error
}

type carService struct {
	carRepo repository.CarRepository
}

func NewCarService(carRepo repository.CarRepository) CarService {
	return &carService{carRepo: carRepo}
}

func (s *carService) CreateCar(ctx context.Context, input car.CreateRequest) (*car.Response, error) {
	if err := validateCreateRequest(input); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	newCar := &car.Car{
		ID:          uuid.New(),
		Model:       input.Model,
		Brand:       input.Brand,
		Year:        input.Year,
		PricePerDay: input.PricePerDay,
		IsAvailable: true,
	}

	if err := s.carRepo.Create(ctx, newCar); err != nil {
		return nil, fmt.Errorf("failed to create car in repository: %w", err)
	}

	return s.entityToResponse(newCar), nil
}

func (s *carService) GetCar(ctx context.Context, id uuid.UUID) (*car.Response, error) {
	carEntity, err := s.carRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get car from repository: %w", err)
	}

	return s.entityToResponse(carEntity), nil
}

func (s *carService) GetAllCars(ctx context.Context, filters map[string]interface{}) ([]car.Response, error) {
	normalizedFilters := normalizeFilters(filters)

	cars, err := s.carRepo.GetAll(ctx, normalizedFilters)
	if err != nil {
		return nil, fmt.Errorf("failed to get cars from repository: %w", err)
	}

	responses := make([]car.Response, 0, len(cars))
	for _, c := range cars {
		responses = append(responses, *s.entityToResponse(c))
	}

	return responses, nil
}

func (s *carService) UpdateCar(ctx context.Context, id uuid.UUID, input car.UpdateRequest) (*car.Response, error) {
	if err := validateUpdateRequest(input); err != nil {
		return nil, fmt.Errorf("validation failed: %w", err)
	}

	if err := s.carRepo.Update(ctx, id, &input); err != nil {
		return nil, fmt.Errorf("failed to update car: %w", err)
	}

	updatedCar, err := s.carRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get updated car: %w", err)
	}

	return s.entityToResponse(updatedCar), nil
}

func (s *carService) DeleteCar(ctx context.Context, id uuid.UUID) error {
	if err := s.carRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete car: %w", err)
	}
	return nil
}

func validatePrice(price float64) error {
	if price <= 0 {
		return fmt.Errorf("price per day must be positive, got: %.2f", price)
	}
	return nil
}

func validateYear(year int) error {
	if year < 1900 || year > 2025 {
		return fmt.Errorf("year must be between 1900 and 2025, got: %d", year)
	}
	return nil
}

func validateModel(model string) error {
	if len(model) == 0 {
		return fmt.Errorf("model cannot be empty")
	}
	return nil
}

func validateBrand(brand string) error {
	if len(brand) == 0 {
		return fmt.Errorf("brand cannot be empty")
	}
	return nil
}

func validateCreateRequest(input car.CreateRequest) error {
	if err := validatePrice(input.PricePerDay); err != nil {
		return err
	}
	if err := validateYear(input.Year); err != nil {
		return err
	}
	if err := validateModel(input.Model); err != nil {
		return err
	}
	if err := validateBrand(input.Brand); err != nil {
		return err
	}
	return nil
}

func validateUpdateRequest(input car.UpdateRequest) error {
	if input.PricePerDay != nil {
		if err := validatePrice(*input.PricePerDay); err != nil {
			return err
		}
	}
	if input.Year != nil {
		if err := validateYear(*input.Year); err != nil {
			return err
		}
	}
	if input.Model != nil {
		if err := validateModel(*input.Model); err != nil {
			return err
		}
	}
	if input.Brand != nil {
		if err := validateBrand(*input.Brand); err != nil {
			return err
		}
	}
	return nil
}

func normalizeFilters(filters map[string]interface{}) map[string]interface{} {
	normalized := make(map[string]interface{})

	for key, value := range filters {
		switch key {
		case "brand":
			if brand, ok := value.(string); ok && brand != "" {
				normalized["brand"] = brand
			}
		case "min_price":
			if minPrice, ok := value.(float64); ok && minPrice >= 0 {
				normalized["min_price"] = minPrice
			}
		case "max_price":
			if maxPrice, ok := value.(float64); ok && maxPrice >= 0 {
				normalized["max_price"] = maxPrice
			}
		case "available":
			if available, ok := value.(bool); ok {
				normalized["available"] = available
			}
		}
	}

	return normalized
}

func (s *carService) entityToResponse(c *car.Car) *car.Response {
	return &car.Response{
		ID:          c.ID.String(),
		Model:       c.Model,
		Brand:       c.Brand,
		Year:        c.Year,
		PricePerDay: c.PricePerDay,
		IsAvailable: c.IsAvailable,
	}
}
