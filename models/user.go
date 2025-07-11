package models

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// User representa el modelo de usuario en la base de datos
type User struct {
	ID        primitive.ObjectID `json:"id" bson:"_id,omitempty" example:"507f1f77bcf86cd799439011"`
	UUID      string            `json:"uuid" bson:"uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string            `json:"name" bson:"name" binding:"required" example:"John Doe"`
	Email     string            `json:"email" bson:"email" binding:"required,email" example:"john.doe@example.com"`
	Age       int               `json:"age" bson:"age" binding:"required,min=1,max=120" example:"30"`
	Phone     string            `json:"phone" bson:"phone" example:"+1234567890"`
	Address   string            `json:"address" bson:"address" example:"123 Main St, City, Country"`
	CreatedAt time.Time         `json:"created_at" bson:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time         `json:"updated_at" bson:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// CreateUserRequest representa la estructura para crear un usuario
type CreateUserRequest struct {
	Name    string `json:"name" binding:"required" example:"John Doe"`
	Email   string `json:"email" binding:"required,email" example:"john.doe@example.com"`
	Age     int    `json:"age" binding:"required,min=1,max=120" example:"30"`
	Phone   string `json:"phone" example:"+1234567890"`
	Address string `json:"address" example:"123 Main St, City, Country"`
}

// UpdateUserRequest representa la estructura para actualizar un usuario
type UpdateUserRequest struct {
	Name    string `json:"name" example:"John Doe"`
	Email   string `json:"email" binding:"omitempty,email" example:"john.doe@example.com"`
	Age     int    `json:"age" binding:"omitempty,min=1,max=120" example:"30"`
	Phone   string `json:"phone" example:"+1234567890"`
	Address string `json:"address" example:"123 Main St, City, Country"`
}

// UserResponse representa la respuesta de usuario
type UserResponse struct {
	ID        string    `json:"id" example:"507f1f77bcf86cd799439011"`
	UUID      string    `json:"uuid" example:"550e8400-e29b-41d4-a716-446655440000"`
	Name      string    `json:"name" example:"John Doe"`
	Email     string    `json:"email" example:"john.doe@example.com"`
	Age       int       `json:"age" example:"30"`
	Phone     string    `json:"phone" example:"+1234567890"`
	Address   string    `json:"address" example:"123 Main St, City, Country"`
	CreatedAt time.Time `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time `json:"updated_at" example:"2023-01-01T00:00:00Z"`
}

// UsersResponse representa la respuesta de lista de usuarios
type UsersResponse struct {
	Users []UserResponse `json:"users"`
	Total int64          `json:"total" example:"10"`
}

// ErrorResponse representa la estructura de respuesta de error
type ErrorResponse struct {
	Error   string `json:"error" example:"Error message"`
	Message string `json:"message" example:"Detailed error message"`
	Code    int    `json:"code" example:"400"`
}

// SuccessResponse representa la estructura de respuesta exitosa
type SuccessResponse struct {
	Message string      `json:"message" example:"Operation completed successfully"`
	Data    interface{} `json:"data,omitempty"`
}

// NewUser crea una nueva instancia de usuario
func NewUser(req CreateUserRequest) *User {
	now := time.Now()
	return &User{
		UUID:      uuid.New().String(),
		Name:      req.Name,
		Email:     req.Email,
		Age:       req.Age,
		Phone:     req.Phone,
		Address:   req.Address,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// ToResponse convierte un User a UserResponse
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID.Hex(),
		UUID:      u.UUID,
		Name:      u.Name,
		Email:     u.Email,
		Age:       u.Age,
		Phone:     u.Phone,
		Address:   u.Address,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// Update actualiza los campos del usuario
func (u *User) Update(req UpdateUserRequest) {
	if req.Name != "" {
		u.Name = req.Name
	}
	if req.Email != "" {
		u.Email = req.Email
	}
	if req.Age > 0 {
		u.Age = req.Age
	}
	if req.Phone != "" {
		u.Phone = req.Phone
	}
	if req.Address != "" {
		u.Address = req.Address
	}
	u.UpdatedAt = time.Now()
} 

