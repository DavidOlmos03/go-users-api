// Package main API Users BRM
//
// API RESTful para gestión de usuarios
//
//	Schemes: http, https
//	Host: localhost:8080
//	BasePath: /api/v1
//	Version: 1.0.0
//
//	Consumes:
//	- application/json
//
//	Produces:
//	- application/json
//
//	Security:
//	- bearer
//
// swagger:meta
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"

	"go-users-api/config"
	"go-users-api/controllers"
	"go-users-api/middleware"
	"go-users-api/repository"
	"go-users-api/services"
	_ "go-users-api/docs"
)

// @title API Users BRM
// @version 1.0
// @description API RESTful para gestión de usuarios con MongoDB
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	// Cargar variables de entorno
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Configurar el modo de Gin
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Inicializar configuración
	cfg := config.NewConfig()

	// Conectar a MongoDB
	client, db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer func() {
		if err := client.Disconnect(context.Background()); err != nil {
			log.Fatal("Error disconnecting from database:", err)
		}
	}()

	// Inicializar repositorios
	userRepo := repository.NewUserRepository(db)

	// Inicializar servicios
	userService := services.NewUserService(userRepo)

	// Inicializar controladores
	userController := controllers.NewUserController(userService)

	// Configurar router
	router := gin.Default()

	// Middleware
	router.Use(middleware.CORS())
	router.Use(middleware.Logger())
	router.Use(middleware.Recovery())

	// Rutas de la API
	api := router.Group("/api/v1")
	{
		// Rutas de usuarios
		users := api.Group("/users")
		{
			users.POST("/", userController.CreateUser)
			users.GET("/", userController.GetUsers)
			users.GET("/:id", userController.GetUserByID)
			users.PUT("/:id", userController.UpdateUser)
			users.DELETE("/:id", userController.DeleteUser)
		}

		// Health check
		api.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"status":  "ok",
				"message": "API Users BRM is running",
				"time":    time.Now().Format(time.RFC3339),
			})
		})
	}

	// Swagger documentation
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Configurar servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	// Canal para manejar señales de terminación
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar servidor en goroutine
	go func() {
		log.Printf("Server starting on port %s", port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error starting server: %v", err)
		}
	}()

	// Esperar señal de terminación
	<-quit
	log.Println("Shutting down server...")

	// Contexto con timeout para shutdown graceful
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exited")
} 

