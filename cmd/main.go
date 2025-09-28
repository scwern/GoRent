package main

import (
	_ "GoRent/docs"
	"GoRent/internal/api"
	"GoRent/internal/config"
	"GoRent/internal/repository"
	"GoRent/internal/repository/analytics"
	"GoRent/internal/repository/db"
	"GoRent/internal/service"
	"GoRent/pkg/jwt"
	"log"
)

// @title           GoRent API
// @version         1.0
// @description     API системы аренды автомобилей
// @host            localhost:8080
// @BasePath        /
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
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
	analyticsRepo := analytics.NewRepository(database)

	jwtManager := jwt.NewManager(cfg.JWT.Secret, cfg.JWT.ExpiresIn)

	authService := service.NewAuthService(userRepo, jwtManager)
	adminService := service.NewAdminService(userRepo)
	carService := service.NewCarService(carRepo)
	rentalService := service.NewRentalService(rentalRepo, carRepo)
	analyticsService := service.NewAnalyticsService(analyticsRepo)

	router := api.SetupRouter(
		authService,
		adminService,
		carService,
		rentalService,
		analyticsService,
		jwtManager,
	)

	if err := router.Run(":8080"); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
