package routes

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"go-users-api/controllers"
	"go-users-api/middleware"
)

// SetupRoutes configura todas las rutas de la aplicación
func SetupRoutes(router *gin.Engine, userController *controllers.UserController) {
	// Middleware global
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())

	// API v1 routes
	api := router.Group("/api/v1")
	{
		// Health check
		api.GET("/health", healthCheck)

		// User routes
		users := api.Group("/users")
		{
			users.POST("", userController.CreateUser)
			users.GET("", userController.GetUsers)
			users.GET("/:id", userController.GetUserByID)
			users.PUT("/:id", userController.UpdateUser)
			users.DELETE("/:id", userController.DeleteUser)
		}
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}

// healthCheck maneja el endpoint de health check
func healthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"message": "API Users BRM is running",
		"time":    time.Now().Format(time.RFC3339),
	})
}
