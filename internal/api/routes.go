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
	jwtManager jwt.Manager,
) *gin.Engine {
	router := gin.Default()

	authHandler := handlers.NewAuthHandler(authService)
	adminHandler := handlers.NewAdminHandler(adminService)

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

	adminGroup := router.Group("/admin")
	adminGroup.Use(authMiddleware, middleware.RoleMiddleware("admin"))
	{
		adminGroup.GET("/users", adminHandler.ListUsers)
		adminGroup.PUT("/users/:id/role", adminHandler.ChangeUserRole)
	}

	return router
}
