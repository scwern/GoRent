package service

import (
	"GoRent/internal/domain/rental"
	"GoRent/internal/repository"
	"GoRent/internal/tests/utils"
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRentalService_CreateRental_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	testDB, cleanup, err := utils.SetupTestDB()
	require.NoError(t, err)
	defer cleanup()

	rentalRepo := repository.NewRentalRepository(testDB)
	carRepo := repository.NewCarRepository(testDB)
	rentalService := NewRentalService(rentalRepo, carRepo)

	carID := uuid.New()
	userID := uuid.New()

	err = utils.CreateTestCar(testDB, carID.String(), "Model S", "Tesla", 2023, 100.0, true)
	require.NoError(t, err)

	err = utils.CreateTestUser(testDB, userID.String(), "test@example.com", "client")
	require.NoError(t, err)

	startDate := time.Now().AddDate(0, 0, 1)
	endDate := startDate.AddDate(0, 0, 3)

	req := rental.CreateRequest{
		CarID:     carID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	result, err := rentalService.CreateRental(context.Background(), req, userID)

	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, carID.String(), result.CarID)
	assert.Equal(t, userID.String(), result.UserID)
	assert.Equal(t, 300.0, result.TotalPrice)
	assert.Equal(t, rental.StatusPending, result.Status)

	rentals, err := rentalRepo.GetByUserID(context.Background(), userID)
	require.NoError(t, err)
	assert.Len(t, rentals, 1)
	assert.Equal(t, result.ID, rentals[0].ID.String())
}

func TestRentalService_CreateRental_CarNotAvailable_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	testDB, cleanup, err := utils.SetupTestDB()
	require.NoError(t, err)
	defer cleanup()

	rentalRepo := repository.NewRentalRepository(testDB)
	carRepo := repository.NewCarRepository(testDB)
	rentalService := NewRentalService(rentalRepo, carRepo)

	carID := uuid.New()
	userID := uuid.New()

	err = utils.CreateTestCar(testDB, carID.String(), "X5", "BMW", 2023, 150.0, false)
	require.NoError(t, err)

	err = utils.CreateTestUser(testDB, userID.String(), "test@example.com", "client")
	require.NoError(t, err)

	req := rental.CreateRequest{
		CarID:     carID,
		StartDate: time.Now().AddDate(0, 0, 1),
		EndDate:   time.Now().AddDate(0, 0, 3),
	}

	result, err := rentalService.CreateRental(context.Background(), req, userID)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "not available")
}

func TestRentalService_GetUserRentals_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	testDB, cleanup, err := utils.SetupTestDB()
	require.NoError(t, err)
	defer cleanup()

	rentalRepo := repository.NewRentalRepository(testDB)
	carRepo := repository.NewCarRepository(testDB)
	rentalService := NewRentalService(rentalRepo, carRepo)

	carID := uuid.New()
	userID := uuid.New()

	err = utils.CreateTestCar(testDB, carID.String(), "Camry", "Toyota", 2022, 50.0, true)
	require.NoError(t, err)

	err = utils.CreateTestUser(testDB, userID.String(), "test@example.com", "client")
	require.NoError(t, err)

	rentals, err := rentalService.GetUserRentals(context.Background(), userID)

	require.NoError(t, err)
	assert.NotNil(t, rentals)
	assert.Len(t, rentals, 0)
}

func TestRentalService_CancelRental_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	testDB, cleanup, err := utils.SetupTestDB()
	require.NoError(t, err)
	defer cleanup()

	rentalRepo := repository.NewRentalRepository(testDB)
	carRepo := repository.NewCarRepository(testDB)
	rentalService := NewRentalService(rentalRepo, carRepo)

	carID := uuid.New()
	userID := uuid.New()

	err = utils.CreateTestCar(testDB, carID.String(), "Model S", "Tesla", 2023, 100.0, true)
	require.NoError(t, err)

	err = utils.CreateTestUser(testDB, userID.String(), "test@example.com", "client")
	require.NoError(t, err)

	startDate := time.Now().AddDate(0, 0, 1)
	endDate := startDate.AddDate(0, 0, 3)

	req := rental.CreateRequest{
		CarID:     carID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	createdRental, err := rentalService.CreateRental(context.Background(), req, userID)
	require.NoError(t, err)

	err = rentalService.CancelRental(context.Background(), uuid.MustParse(createdRental.ID), userID)

	require.NoError(t, err)

	updatedRental, err := rentalRepo.GetByID(context.Background(), uuid.MustParse(createdRental.ID))
	require.NoError(t, err)
	assert.Equal(t, rental.StatusCanceled, updatedRental.Status)
}

func TestRentalService_ConcurrentRentals_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	testDB, cleanup, err := utils.SetupTestDB()
	require.NoError(t, err)
	defer cleanup()

	rentalRepo := repository.NewRentalRepository(testDB)
	carRepo := repository.NewCarRepository(testDB)

	carID := uuid.New()
	userID1 := uuid.New()
	userID2 := uuid.New()

	err = utils.CreateTestCar(testDB, carID.String(), "Model S", "Tesla", 2023, 100.0, true)
	require.NoError(t, err)

	err = utils.CreateTestUser(testDB, userID1.String(), "user1@example.com", "client")
	require.NoError(t, err)

	err = utils.CreateTestUser(testDB, userID2.String(), "user2@example.com", "client")
	require.NoError(t, err)

	startDate := time.Now().AddDate(0, 0, 1)
	endDate := startDate.AddDate(0, 0, 2)
	req := rental.CreateRequest{
		CarID:     carID,
		StartDate: startDate,
		EndDate:   endDate,
	}

	results := make(chan struct {
		rental *rental.Response
		err    error
	}, 2)

	for i := 0; i < 2; i++ {
		go func(userID uuid.UUID) {
			rentalService := NewRentalService(rentalRepo, carRepo)
			result, err := rentalService.CreateRental(context.Background(), req, userID)
			results <- struct {
				rental *rental.Response
				err    error
			}{result, err}
		}([]uuid.UUID{userID1, userID2}[i])
	}

	var successes int
	var errors int

	for i := 0; i < 2; i++ {
		result := <-results
		if result.err != nil {
			errors++
			t.Logf("Error (expected for one request): %v", result.err)
		} else {
			successes++
			assert.NotNil(t, result.rental)
		}
	}

	assert.Equal(t, 1, successes)
	assert.Equal(t, 1, errors)
}
