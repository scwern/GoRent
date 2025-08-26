package api

import (
	"GoRent/internal/api/handlers"
	"GoRent/internal/api/middleware"
	"GoRent/internal/service"
	"GoRent/pkg/jwt"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authService service.AuthService,
	adminService service.AdminService,
	carService service.CarService,
	rentalService service.RentalService,
	jwtManager jwt.Manager,
) *gin.Engine {
	router := gin.Default()

	authHandler := handlers.NewAuthHandler(authService)
	adminHandler := handlers.NewAdminHandler(adminService)
	carHandler := handlers.NewCarHandler(carService)
	rentalHandler := handlers.NewRentalHandler(rentalService)

	authMiddleware := middleware.AuthMiddleware(jwtManager)

	authGroup := router.Group("/auth")
	{
		authGroup.POST("/register", authHandler.Register)
		authGroup.POST("/login", authHandler.Login)

		authGroup.GET("/validate", authMiddleware, func(c *gin.Context) {
			userID, _ := c.Get("userID")
			userRole, _ := c.Get("userRole")

			c.JSON(200, gin.H{
				"status":  "success",
				"user_id": userID,
				"role":    userRole,
			})
		})
	}

	carGroup := router.Group("/cars")
	{
		carGroup.GET("", carHandler.GetAllCars)
		carGroup.GET("/:id", carHandler.GetCar)
	}

	protectedCarGroup := router.Group("/cars")
	protectedCarGroup.Use(authMiddleware, middleware.RoleMiddleware("manager", "admin"))
	{
		protectedCarGroup.POST("", carHandler.CreateCar)
		protectedCarGroup.PUT("/:id", carHandler.UpdateCar)
		protectedCarGroup.DELETE("/:id", carHandler.DeleteCar)
	}

	rentalGroup := router.Group("/rentals")
	rentalGroup.Use(authMiddleware)
	{
		rentalGroup.POST("", rentalHandler.CreateRental)
		rentalGroup.GET("", rentalHandler.GetUserRentals)
		rentalGroup.GET("/:id", rentalHandler.GetRental)
		rentalGroup.PUT("/:id/cancel", rentalHandler.CancelRental)
	}

	managerRentalGroup := router.Group("/rentals")
	managerRentalGroup.Use(authMiddleware, middleware.RoleMiddleware("manager", "admin"))
	{
		managerRentalGroup.PUT("/:id/approve", rentalHandler.ApproveRental)
	}

	adminGroup := router.Group("/admin")
	adminGroup.Use(authMiddleware, middleware.RoleMiddleware("admin"))
	{
		adminGroup.GET("/users", adminHandler.ListUsers)
		adminGroup.PUT("/users/:id/role", adminHandler.ChangeUserRole)
	}

	return router
}
