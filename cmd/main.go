package main

import (
	"GoRent/internal/api"
	"GoRent/internal/config"
	"GoRent/internal/repository"
	"GoRent/internal/repository/db"
	"GoRent/internal/service"
	"GoRent/pkg/jwt"
	"log"
)

func main() {
	cfg := config.Load()

	database, err := db.NewPostgresDB(cfg)
	if err != nil {
		log.Fatalf("failed to initialize database: %v", err)
	}
	defer database.Close()

	userRepo := repository.NewUserRepository(database)
	carRepo := repository.NewCarRepository(database)
	rentalRepo := repository.NewRentalRepository(database)

	jwtManager := jwt.NewManager(cfg.JWT.Secret, cfg.JWT.ExpiresIn)

	authService := service.NewAuthService(userRepo, jwtManager)
	adminService := service.NewAdminService(userRepo)
	carService := service.NewCarService(carRepo)
	rentalService := service.NewRentalService(rentalRepo, carRepo)

	router := api.SetupRouter(authService, adminService, carService, rentalService, jwtManager)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
