package config

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config estructura para manejar la configuración de la aplicación
type Config struct {
	MongoURI      string
	MongoDatabase string
	Port          string
	GinMode       string
}

// NewConfig crea una nueva instancia de configuración
func NewConfig() *Config {
	return &Config{
		MongoURI:      getEnv("MONGO_URI", "mongodb://localhost:27017"),
		MongoDatabase: getEnv("MONGO_DATABASE", "users_brm"),
		Port:          getEnv("PORT", "8080"),
		GinMode:       getEnv("GIN_MODE", "debug"),
	}
}

// getEnv obtiene una variable de entorno o retorna un valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// ConnectDB establece la conexión con MongoDB
func ConnectDB(cfg *Config) (*mongo.Client, *mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(cfg.MongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, nil, err
	}

	// Verificar la conexión
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	log.Println("Connected to MongoDB successfully")
	return client, client.Database(cfg.MongoDatabase), nil
} 

