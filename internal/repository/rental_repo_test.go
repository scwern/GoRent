package repository

import (
	"context"
	"testing"
	"time"

	"GoRent/internal/domain/rental"
	"GoRent/internal/tests/utils"
	"github.com/google/uuid"
)

func TestRentalRepository_RaceCondition(t *testing.T) {
	testDB, cleanup, err := utils.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test DB: %v", err)
	}
	defer cleanup()

	repo := NewRentalRepository(testDB)

	carID, _ := uuid.Parse("44444444-4444-4444-4444-444444444444")
	userID, _ := uuid.Parse("33333333-3333-3333-3333-333333333333")

	startDate := time.Now().AddDate(0, 0, 1)
	endDate := startDate.AddDate(0, 0, 3)

	results := make(chan error, 2)

	for i := 0; i < 2; i++ {
		go func() {
			available, err := repo.CheckCarAvailability(context.Background(), carID, startDate, endDate)
			if err != nil {
				results <- err
				return
			}
			if available {
				testRental := &rental.Rental{
					ID:         uuid.New(),
					CarID:      carID,
					UserID:     userID,
					StartDate:  startDate,
					EndDate:    endDate,
					TotalPrice: 300.0,
					Status:     rental.StatusPending,
					CreatedAt:  time.Now(),
				}
				err := repo.Create(context.Background(), testRental)
				results <- err
			} else {
				results <- nil
			}
		}()
	}

	var successes int
	for i := 0; i < 2; i++ {
		err := <-results
		if err == nil {
			successes++
		} else {
			t.Logf("Error: %v", err)
		}
	}

	if successes != 1 {
		t.Errorf("Expected 1 success, got %d", successes)
	}
}
