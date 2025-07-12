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

	"go-users-api/config"
	"go-users-api/controllers"
	_ "go-users-api/docs"
	"go-users-api/repository"
	"go-users-api/routes"
	"go-users-api/services"
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
	// Cargar variables de entorno desde .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Inicializar configuración
	cfg := config.NewConfig()

	// Configurar el modo de Gin desde la configuración
	if cfg.GinMode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

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

	// Configurar rutas
	routes.SetupRoutes(router, userController)

	// Configurar servidor usando la configuración
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	// Canal para manejar señales de terminación
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Iniciar servidor en goroutine
	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
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
